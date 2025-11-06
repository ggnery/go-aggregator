BEGIN;

-- Keep objects out of 'public'
CREATE SCHEMA IF NOT EXISTS orchestrator AUTHORIZATION CURRENT_USER;
SET search_path = orchestrator, public;

-- 1) Accounts (owning identities)
CREATE TABLE IF NOT EXISTS account (
  account_id     uuid PRIMARY KEY,
  wallet_address text,
  created_at     timestamptz NOT NULL DEFAULT now()
);

-- 2) Models & artifacts (base first, then dependents)
CREATE TABLE IF NOT EXISTS model (
  model_id            uuid PRIMARY KEY,
  family              text NOT NULL,           -- "llama3", "whisper", ...
  framework           text NOT NULL,           -- "pytorch","onnx","tflite"
  created_at          timestamptz NOT NULL DEFAULT now(),
  require_accelerator boolean NOT NULL DEFAULT true,
  meta                jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS artifact (
  artifact_id       uuid PRIMARY KEY,
  family            text NOT NULL,    -- "model","code","dataset","tokenizer", ...
  resource_locator  text NOT NULL,    -- URI/URL
  integrity_hash    text,
  required_scope    text,             -- e.g., "netflix","hbo"
  created_at        timestamptz NOT NULL DEFAULT now(),
  annotations       jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS artifact_version (
  artifact_version_id uuid PRIMARY KEY,
  artifact_id         uuid NOT NULL REFERENCES artifact(artifact_id) ON DELETE CASCADE,
  created_at          timestamptz NOT NULL DEFAULT now(),
  meta                jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS model_variant (
  variant_id         uuid PRIMARY KEY,
  model_id           uuid NOT NULL REFERENCES model(model_id) ON DELETE CASCADE,
  created_at         timestamptz NOT NULL DEFAULT now(),
  precision          text,          -- "fp16","int8","fp8","q4_k_m", ...
  device_class       text,          -- "gpu","tpu","ane","cpu"
  min_vram_mb        integer,
  artifact_version_id uuid REFERENCES artifact_version(artifact_version_id),
  UNIQUE (model_id, precision, device_class, artifact_version_id)
);
CREATE INDEX IF NOT EXISTS model_variant_device_idx
  ON model_variant(device_class, min_vram_mb);

-- 3) Nodes & capabilities
CREATE TABLE IF NOT EXISTS node (
  node_id              uuid PRIMARY KEY,
  enrolled_at          timestamptz NOT NULL DEFAULT now(),
  agent_version        text NOT NULL,
  node_type            text NOT NULL,
  owner_account_id     uuid REFERENCES account(account_id),
  operating_account_id uuid REFERENCES account(account_id),
  description          jsonb NOT NULL, -- {cpu_cores, gpu_model, vram_mb, ram_mb, disk_mb}
  region               text
);
CREATE INDEX IF NOT EXISTS node_by_type_owner_idx
  ON node(node_type, owner_account_id, operating_account_id);

-- Flexible capability surface for scheduling filters
CREATE TABLE IF NOT EXISTS node_capability (
  node_id  uuid PRIMARY KEY REFERENCES node(node_id) ON DELETE CASCADE,
  cap      jsonb NOT NULL  -- e.g. {"framework":"torch","cuda":12.1,"max_batch":8,"supports_fp8":true}
);
CREATE INDEX IF NOT EXISTS node_cap_gin ON node_capability USING GIN (cap jsonb_path_ops);

CREATE TABLE IF NOT EXISTS node_session (
  session_id        uuid PRIMARY KEY,
  node_id           uuid NOT NULL REFERENCES node(node_id) ON DELETE CASCADE,
  started_at        timestamptz NOT NULL DEFAULT now(),
  ended_at          timestamptz,
  status            text NOT NULL CHECK (status IN ('active','draining','lost','ended')),
  controllable      boolean NOT NULL DEFAULT true,
  last_heartbeat_at timestamptz,
  rtt_ms            integer,
  health_score      real,
  metrics           jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS node_session_node_idx
  ON node_session(node_id, status);

CREATE TABLE IF NOT EXISTS node_cache_entry (
  node_id             uuid NOT NULL REFERENCES node(node_id) ON DELETE CASCADE,
  artifact_version_id uuid NOT NULL REFERENCES artifact_version(artifact_version_id),
  acquired_at         timestamptz NOT NULL DEFAULT now(),
  last_used_at        timestamptz,
  size_bytes          bigint,
  PRIMARY KEY (node_id, artifact_version_id)
);
CREATE INDEX IF NOT EXISTS node_cache_artifact_idx
  ON node_cache_entry(artifact_version_id);

-- 4) Jobs
CREATE TABLE IF NOT EXISTS job (
  job_id          uuid PRIMARY KEY,
  owner_id        uuid NOT NULL REFERENCES account(account_id),
  project_id      uuid,
  priority        integer NOT NULL DEFAULT 0,
  sla_tier        text NOT NULL DEFAULT 'standard', -- 'latency','throughput','best_effort'
  status          text NOT NULL CHECK (status IN ('queued','running','aggregating','completed','failed','cancelled')),
  submitted_at    timestamptz NOT NULL DEFAULT now(),
  started_at      timestamptz,
  finished_at     timestamptz,
  output_ref      text,
  parent_job      uuid REFERENCES job(job_id),
  idempotency_key text,
  code_reference  text,
  gas_fee_est     numeric(38,18),
  gas_fee_actual  numeric(38,18),
  deadline_at     timestamptz,
  annotations     jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS job_status_prio_submitted_idx
  ON job(status, priority DESC, submitted_at);
-- Deduplicate per owner when idempotency_key is provided
CREATE UNIQUE INDEX IF NOT EXISTS job_owner_idem_uidx
  ON job(owner_id, idempotency_key)
  WHERE idempotency_key IS NOT NULL;

-- 5) Tasks
CREATE TABLE IF NOT EXISTS task (
  task_id        uuid PRIMARY KEY,
  parent_job     uuid NOT NULL REFERENCES job(job_id) ON DELETE CASCADE,
  status         text NOT NULL CHECK (status IN ('pending','leased','completed','failed','timeout','cancelled','quarantined')),
  created_at     timestamptz NOT NULL DEFAULT now(),
  ready_at       timestamptz NOT NULL DEFAULT now(), -- backoff / not-before time
  priority       integer NOT NULL DEFAULT 0,
  attempts_made  integer NOT NULL DEFAULT 0,
  max_attempts   integer NOT NULL DEFAULT 3,
  models_needed  jsonb NOT NULL DEFAULT '[]'::jsonb,  -- [{family, min_mem, precision, ...}]
  input_ref      jsonb,
  code_reference text,
  deadline_at    timestamptz,
  annotations    jsonb NOT NULL DEFAULT '{}'::jsonb,
  CHECK (attempts_made <= max_attempts)
);
-- Indexing strategy for runnable tasks:
CREATE INDEX IF NOT EXISTS task_pending_ready_prio_idx
  ON task (ready_at, priority DESC, created_at)
  WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS task_parent_idx ON task(parent_job);
CREATE INDEX IF NOT EXISTS task_status_idx ON task(status);

-- 6) Attempts (leases live here)
CREATE TABLE IF NOT EXISTS task_attempt (
  attempt_id        bigserial PRIMARY KEY,
  task_id           uuid NOT NULL REFERENCES task(task_id) ON DELETE CASCADE,
  attempt_no        integer NOT NULL, -- 1..N
  session_id        uuid REFERENCES node_session(session_id),
  node_id           uuid REFERENCES node(node_id),
  lease_token       uuid NOT NULL,
  lease_expires_at  timestamptz NOT NULL,
  leased_at         timestamptz NOT NULL DEFAULT now(),
  status            text NOT NULL CHECK (status IN ('leased','running','succeeded','failed','aborted','expired')),
  started_at        timestamptz,
  finished_at       timestamptz,
  result_ref        text,
  error_code        text,
  error_message     text,
  metrics           jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE UNIQUE INDEX IF NOT EXISTS task_attempt_unique_no
  ON task_attempt(task_id, attempt_no);
-- For "what's the active lease for this task?"
CREATE INDEX IF NOT EXISTS task_attempt_active_by_task_idx
  ON task_attempt (task_id, lease_expires_at)
  WHERE status IN ('leased','running');

-- For the reaper sweep: "which leases have expired?"
CREATE INDEX IF NOT EXISTS task_attempt_expiry_idx
  ON task_attempt (lease_expires_at)
  WHERE status IN ('leased','running');

-- 7) Aggregation
CREATE TABLE IF NOT EXISTS aggregation (
  job_id             uuid PRIMARY KEY REFERENCES job(job_id) ON DELETE CASCADE,
  expected_parts     integer NOT NULL,
  received_parts     integer NOT NULL DEFAULT 0,
  strategy           text NOT NULL, -- 'concat','reduce','map-reduce','majority','custom'
  merge_status       text NOT NULL CHECK (merge_status IN ('idle','merging','finalized','failed')),
  final_ref          text,
  last_part_merged_at timestamptz,
  annotations        jsonb NOT NULL DEFAULT '{}'::jsonb,
  CHECK (received_parts <= expected_parts)
);
CREATE INDEX IF NOT EXISTS aggregation_merge_status_idx ON aggregation(merge_status);

CREATE TABLE IF NOT EXISTS aggregation_part (
  job_id       uuid NOT NULL REFERENCES aggregation(job_id) ON DELETE CASCADE,
  part_key     text NOT NULL,  -- stable dedupe key: (partition_index || chunk_id || checksum)
  task_id      uuid REFERENCES task(task_id),
  attempt_id   bigint REFERENCES task_attempt(attempt_id),
  received_at  timestamptz NOT NULL DEFAULT now(),
  payload_ref  text NOT NULL,
  merge_state  text NOT NULL CHECK (merge_state IN ('pending','merged','skipped','corrupt')),
  PRIMARY KEY (job_id, part_key)
);

-- 8) Outbox (CDC)
CREATE TABLE IF NOT EXISTS outbox_event (
  outbox_id    bigserial PRIMARY KEY,
  created_at   timestamptz NOT NULL DEFAULT now(),
  topic        text NOT NULL,   -- "metrics","logs","webhook"
  payload      jsonb NOT NULL,
  delivered_at timestamptz
);
CREATE INDEX IF NOT EXISTS outbox_undelivered_idx
  ON outbox_event(created_at) WHERE delivered_at IS NULL;

-- 9) Canonical events (partitioned by month) — PK must include partition key `ts`
CREATE TABLE IF NOT EXISTS event (
  event_id  bigint GENERATED ALWAYS AS IDENTITY,
  ts        timestamptz NOT NULL DEFAULT now(),
  entity_type text NOT NULL,
  entity_id   text NOT NULL,
  event_type  text NOT NULL,
  data        jsonb NOT NULL DEFAULT '{}'::jsonb,
  PRIMARY KEY (ts, event_id)          
) PARTITION BY RANGE (ts);

-- Global (partitioned) indexes—Postgres will create per-partition indexes for existing partitions
-- We'll also ensure new partitions get these indexes when created (see block below).
CREATE INDEX IF NOT EXISTS event_entity_idx ON event(entity_type, entity_id, ts DESC);
CREATE INDEX IF NOT EXISTS event_type_ts_idx ON event(event_type, ts DESC);

COMMIT;
