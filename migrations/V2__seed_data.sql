-- Seed data for development/testing
SET search_path = orchestrator, public;

-- ========================================
-- 1) ACCOUNTS
-- ========================================
INSERT INTO account (account_id, wallet_address, created_at) VALUES
  ('00000000-0000-0000-0000-000000000001', '0x1234567890123456789012345678901234567890', now() - interval '30 days'),
  ('00000000-0000-0000-0000-000000000002', '0xabcdefabcdefabcdefabcdefabcdefabcdefabcd', now() - interval '25 days'),
  ('00000000-0000-0000-0000-000000000003', '0x9876543210987654321098765432109876543210', now() - interval '20 days');

-- ========================================
-- 2) MODELS & ARTIFACTS
-- ========================================

-- Models
INSERT INTO model (model_id, family, framework, created_at, require_accelerator, meta) VALUES
  ('10000000-0000-0000-0000-000000000001', 'llama3', 'pytorch', now() - interval '60 days', true, '{"size":"8B","context_length":8192}'::jsonb),
  ('10000000-0000-0000-0000-000000000002', 'llama3', 'pytorch', now() - interval '55 days', true, '{"size":"70B","context_length":8192}'::jsonb),
  ('10000000-0000-0000-0000-000000000003', 'whisper', 'onnx', now() - interval '50 days', false, '{"size":"large-v3","languages":"multilingual"}'::jsonb),
  ('10000000-0000-0000-0000-000000000004', 'stable-diffusion', 'pytorch', now() - interval '45 days', true, '{"version":"xl-1.0","resolution":"1024x1024"}'::jsonb);

-- Artifacts
INSERT INTO artifact (artifact_id, family, resource_locator, integrity_hash, required_scope, created_at, annotations) VALUES
  ('20000000-0000-0000-0000-000000000001', 'model', 'ipfs://QmLlama38BWeights', 'sha256:abc123...', NULL, now() - interval '60 days', '{"size_gb":16}'::jsonb),
  ('20000000-0000-0000-0000-000000000002', 'model', 'ipfs://QmLlama370BWeights', 'sha256:def456...', NULL, now() - interval '55 days', '{"size_gb":140}'::jsonb),
  ('20000000-0000-0000-0000-000000000003', 'model', 'ipfs://QmWhisperLarge', 'sha256:ghi789...', NULL, now() - interval '50 days', '{"size_gb":3}'::jsonb),
  ('20000000-0000-0000-0000-000000000004', 'tokenizer', 'ipfs://QmLlamaTokenizer', 'sha256:jkl012...', NULL, now() - interval '60 days', '{"vocab_size":32000}'::jsonb),
  ('20000000-0000-0000-0000-000000000005', 'dataset', 's3://training-data/common-crawl-sample', 'sha256:mno345...', NULL, now() - interval '70 days', '{"size_gb":500}'::jsonb);

-- Artifact Versions
INSERT INTO artifact_version (artifact_version_id, artifact_id, created_at, meta) VALUES
  ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', now() - interval '60 days', '{"version":"v1.0"}'::jsonb),
  ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000002', now() - interval '55 days', '{"version":"v1.0"}'::jsonb),
  ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000003', now() - interval '50 days', '{"version":"v3.0"}'::jsonb),
  ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000004', now() - interval '60 days', '{"version":"v1.0"}'::jsonb);

-- Model Variants (different precision/device combinations)
INSERT INTO model_variant (variant_id, model_id, created_at, precision, device_class, min_vram_mb, artifact_version_id) VALUES
  -- Llama3 8B variants
  ('40000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', now() - interval '60 days', 'fp16', 'gpu', 16384, '30000000-0000-0000-0000-000000000001'),
  ('40000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', now() - interval '60 days', 'int8', 'gpu', 8192, '30000000-0000-0000-0000-000000000001'),
  ('40000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', now() - interval '60 days', 'q4_k_m', 'cpu', 4096, '30000000-0000-0000-0000-000000000001'),
  
  -- Llama3 70B variants
  ('40000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000002', now() - interval '55 days', 'fp16', 'gpu', 140000, '30000000-0000-0000-0000-000000000002'),
  ('40000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000002', now() - interval '55 days', 'int8', 'gpu', 70000, '30000000-0000-0000-0000-000000000002'),
  
  -- Whisper variants
  ('40000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000003', now() - interval '50 days', 'fp32', 'cpu', 2048, '30000000-0000-0000-0000-000000000003'),
  ('40000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000003', now() - interval '50 days', 'fp16', 'gpu', 1024, '30000000-0000-0000-0000-000000000003');

-- ========================================
-- 3) NODES & CAPABILITIES
-- ========================================

INSERT INTO node (node_id, enrolled_at, agent_version, node_type, owner_account_id, operating_account_id, description, region) VALUES
  -- GPU nodes
  ('50000000-0000-0000-0000-000000000001', now() - interval '30 days', 'v1.2.0', 'compute', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 
   '{"cpu_cores":32,"gpu_model":"RTX 4090","vram_mb":24576,"ram_mb":65536,"disk_mb":1000000}'::jsonb, 'us-west-2'),
  
  ('50000000-0000-0000-0000-000000000002', now() - interval '25 days', 'v1.2.1', 'compute', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 
   '{"cpu_cores":64,"gpu_model":"A100","vram_mb":81920,"ram_mb":131072,"disk_mb":2000000}'::jsonb, 'us-east-1'),
  
  ('50000000-0000-0000-0000-000000000003', now() - interval '20 days', 'v1.2.1', 'compute', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 
   '{"cpu_cores":16,"gpu_model":"RTX 3090","vram_mb":24576,"ram_mb":32768,"disk_mb":500000}'::jsonb, 'eu-west-1'),
  
  -- CPU nodes
  ('50000000-0000-0000-0000-000000000004', now() - interval '15 days', 'v1.2.0', 'compute', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 
   '{"cpu_cores":128,"ram_mb":262144,"disk_mb":2000000}'::jsonb, 'us-west-1'),
  
  ('50000000-0000-0000-0000-000000000005', now() - interval '10 days', 'v1.2.1', 'compute', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000003', 
   '{"cpu_cores":64,"ram_mb":131072,"disk_mb":1000000}'::jsonb, 'ap-southeast-1');

-- Node capabilities
INSERT INTO node_capability (node_id, cap) VALUES
  ('50000000-0000-0000-0000-000000000001', '{"framework":"pytorch","cuda":"12.1","max_batch":8,"supports_fp16":true,"supports_int8":true}'::jsonb),
  ('50000000-0000-0000-0000-000000000002', '{"framework":"pytorch","cuda":"12.2","max_batch":32,"supports_fp16":true,"supports_int8":true,"supports_fp8":true}'::jsonb),
  ('50000000-0000-0000-0000-000000000003', '{"framework":"pytorch","cuda":"11.8","max_batch":4,"supports_fp16":true,"supports_int8":true}'::jsonb),
  ('50000000-0000-0000-0000-000000000004', '{"framework":"onnx","max_batch":16}'::jsonb),
  ('50000000-0000-0000-0000-000000000005', '{"framework":"onnx","max_batch":8}'::jsonb);

-- Node sessions (some active, some ended)
INSERT INTO node_session (session_id, node_id, started_at, ended_at, status, controllable, last_heartbeat_at, rtt_ms, health_score, metrics) VALUES
  ('60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', now() - interval '5 hours', NULL, 'active', true, now() - interval '30 seconds', 25, 0.98, '{"tasks_completed":152}'::jsonb),
  ('60000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', now() - interval '3 hours', NULL, 'active', true, now() - interval '45 seconds', 35, 0.95, '{"tasks_completed":89}'::jsonb),
  ('60000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000003', now() - interval '2 hours', NULL, 'active', true, now() - interval '20 seconds', 60, 0.92, '{"tasks_completed":45}'::jsonb),
  ('60000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000004', now() - interval '6 hours', now() - interval '1 hour', 'ended', true, now() - interval '1 hour', 15, 0.99, '{"tasks_completed":203}'::jsonb),
  ('60000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000005', now() - interval '1 hour', NULL, 'active', true, now() - interval '15 seconds', 120, 0.88, '{"tasks_completed":12}'::jsonb);

-- Node cache entries (models cached on nodes)
INSERT INTO node_cache_entry (node_id, artifact_version_id, acquired_at, last_used_at, size_bytes) VALUES
  ('50000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', now() - interval '4 days', now() - interval '2 hours', 17179869184),
  ('50000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000004', now() - interval '4 days', now() - interval '2 hours', 1073741824),
  ('50000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', now() - interval '3 days', now() - interval '1 hour', 150323855360),
  ('50000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000001', now() - interval '2 days', now() - interval '3 hours', 17179869184),
  ('50000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', now() - interval '5 days', now() - interval '4 hours', 3221225472);

-- ========================================
-- 4) JOBS & TASKS (Sample workload)
-- ========================================

-- Job 1: Completed inference job
INSERT INTO job (job_id, owner_id, project_id, priority, sla_tier, status, submitted_at, started_at, finished_at, output_ref, code_reference, gas_fee_est, gas_fee_actual, annotations) VALUES
  ('70000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', NULL, 10, 'latency', 'completed', 
   now() - interval '10 hours', now() - interval '9 hours 55 minutes', now() - interval '9 hours 30 minutes', 
   's3://results/job-70000000-0000-0000-0000-000000000001.json', 
   'ipfs://QmCodeHash123', 
   0.0015, 0.0012, 
   '{"model":"llama3-8b","task_type":"text_generation","batch_size":10}'::jsonb);

-- Tasks for Job 1 (all completed)
INSERT INTO task (task_id, parent_job, status, created_at, ready_at, priority, attempts_made, max_attempts, models_needed, input_ref, code_reference, annotations) VALUES
  ('80000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', 'completed', now() - interval '10 hours', now() - interval '10 hours', 10, 1, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"prompt":"What is the capital of France?","max_tokens":100}'::jsonb, 
   'ipfs://QmCodeHash123', 
   '{"partition":1}'::jsonb),
  
  ('80000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', 'completed', now() - interval '10 hours', now() - interval '10 hours', 10, 1, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"prompt":"Explain quantum computing","max_tokens":200}'::jsonb, 
   'ipfs://QmCodeHash123', 
   '{"partition":2}'::jsonb);

-- Job 2: Currently running job with aggregation
INSERT INTO job (job_id, owner_id, project_id, priority, sla_tier, status, submitted_at, started_at, code_reference, gas_fee_est, annotations) VALUES
  ('70000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', NULL, 5, 'throughput', 'running', 
   now() - interval '2 hours', now() - interval '1 hour 55 minutes', 
   'ipfs://QmCodeHash456', 
   0.0250, 
   '{"model":"llama3-8b","task_type":"batch_inference","total_items":100}'::jsonb);

-- Aggregation for Job 2
INSERT INTO aggregation (job_id, expected_parts, received_parts, strategy, merge_status, annotations) VALUES
  ('70000000-0000-0000-0000-000000000002', 100, 45, 'concat', 'idle', '{"output_format":"jsonl"}'::jsonb);

-- Tasks for Job 2 (mixed states)
INSERT INTO task (task_id, parent_job, status, created_at, ready_at, priority, attempts_made, max_attempts, models_needed, input_ref, code_reference, annotations) VALUES
  ('80000000-0000-0000-0000-000000000010', '70000000-0000-0000-0000-000000000002', 'completed', now() - interval '2 hours', now() - interval '2 hours', 5, 1, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"batch_start":0,"batch_end":10}'::jsonb, 
   'ipfs://QmCodeHash456', 
   '{"partition":1}'::jsonb),
  
  ('80000000-0000-0000-0000-000000000011', '70000000-0000-0000-0000-000000000002', 'completed', now() - interval '2 hours', now() - interval '2 hours', 5, 1, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"batch_start":10,"batch_end":20}'::jsonb, 
   'ipfs://QmCodeHash456', 
   '{"partition":2}'::jsonb),
  
  ('80000000-0000-0000-0000-000000000012', '70000000-0000-0000-0000-000000000002', 'leased', now() - interval '2 hours', now() - interval '2 hours', 5, 1, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"batch_start":20,"batch_end":30}'::jsonb, 
   'ipfs://QmCodeHash456', 
   '{"partition":3}'::jsonb),
  
  ('80000000-0000-0000-0000-000000000013', '70000000-0000-0000-0000-000000000002', 'pending', now() - interval '2 hours', now() - interval '2 hours', 5, 0, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"batch_start":30,"batch_end":40}'::jsonb, 
   'ipfs://QmCodeHash456', 
   '{"partition":4}'::jsonb),
  
  ('80000000-0000-0000-0000-000000000014', '70000000-0000-0000-0000-000000000002', 'pending', now() - interval '2 hours', now() - interval '2 hours', 5, 0, 3, 
   '[{"family":"llama3","min_vram_mb":8192,"precision":"int8"}]'::jsonb, 
   '{"batch_start":40,"batch_end":50}'::jsonb, 
   'ipfs://QmCodeHash456', 
   '{"partition":5}'::jsonb);

-- Job 3: Queued job (not started yet)
INSERT INTO job (job_id, owner_id, project_id, priority, sla_tier, status, submitted_at, code_reference, gas_fee_est, annotations) VALUES
  ('70000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000003', NULL, 15, 'latency', 'queued', 
   now() - interval '30 minutes', 
   'ipfs://QmCodeHash789', 
   0.0080, 
   '{"model":"whisper","task_type":"transcription"}'::jsonb);

-- ========================================
-- 5) TASK ATTEMPTS
-- ========================================

-- Attempts for completed tasks (explicitly setting attempt_id to match aggregation_part references)
INSERT INTO task_attempt (attempt_id, task_id, attempt_no, session_id, node_id, lease_token, lease_expires_at, leased_at, status, started_at, finished_at, result_ref, metrics) VALUES
  (1, '80000000-0000-0000-0000-000000000001', 1, '60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 
   '90000000-0000-0000-0000-000000000001', now() - interval '9 hours 30 minutes', now() - interval '9 hours 50 minutes', 
   'succeeded', now() - interval '9 hours 45 minutes', now() - interval '9 hours 40 minutes', 
   's3://results/task-80000000-0000-0000-0000-000000000001.json', 
   '{"duration_ms":300000,"tokens_generated":85}'::jsonb),
  
  (2, '80000000-0000-0000-0000-000000000002', 1, '60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 
   '90000000-0000-0000-0000-000000000002', now() - interval '9 hours 35 minutes', now() - interval '9 hours 45 minutes', 
   'succeeded', now() - interval '9 hours 42 minutes', now() - interval '9 hours 38 minutes', 
   's3://results/task-80000000-0000-0000-0000-000000000002.json', 
   '{"duration_ms":240000,"tokens_generated":175}'::jsonb),
  
  (3, '80000000-0000-0000-0000-000000000010', 1, '60000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', 
   '90000000-0000-0000-0000-000000000010', now() - interval '1 hour 30 minutes', now() - interval '1 hour 45 minutes', 
   'succeeded', now() - interval '1 hour 42 minutes', now() - interval '1 hour 35 minutes', 
   's3://results/task-80000000-0000-0000-0000-000000000010.json', 
   '{"duration_ms":420000,"batch_size":10}'::jsonb),
  
  (4, '80000000-0000-0000-0000-000000000011', 1, '60000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000003', 
   '90000000-0000-0000-0000-000000000011', now() - interval '1 hour 20 minutes', now() - interval '1 hour 40 minutes', 
   'succeeded', now() - interval '1 hour 38 minutes', now() - interval '1 hour 25 minutes', 
   's3://results/task-80000000-0000-0000-0000-000000000011.json', 
   '{"duration_ms":780000,"batch_size":10}'::jsonb),
  
  -- Active lease
  (5, '80000000-0000-0000-0000-000000000012', 1, '60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 
   '90000000-0000-0000-0000-000000000012', now() + interval '15 minutes', now() - interval '5 minutes', 
   'running', now() - interval '4 minutes', NULL, NULL, 
   '{"progress":0.6}'::jsonb);

-- Update the sequence to continue from the last inserted ID
SELECT setval('task_attempt_attempt_id_seq', 5, true);

-- Aggregation parts for Job 2
INSERT INTO aggregation_part (job_id, part_key, task_id, attempt_id, received_at, payload_ref, merge_state) VALUES
  ('70000000-0000-0000-0000-000000000002', 'part-001', '80000000-0000-0000-0000-000000000010', 3, now() - interval '1 hour 30 minutes', 's3://results/task-80000000-0000-0000-0000-000000000010.json', 'merged'),
  ('70000000-0000-0000-0000-000000000002', 'part-002', '80000000-0000-0000-0000-000000000011', 4, now() - interval '1 hour 20 minutes', 's3://results/task-80000000-0000-0000-0000-000000000011.json', 'merged');

-- ========================================
-- 6) EVENTS (Sample activity)
-- ========================================

-- Create partition for current month (event table is partitioned by timestamp)
-- This will hold events for November 2025
CREATE TABLE IF NOT EXISTS event_2025_11 PARTITION OF event
  FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');

-- Create indexes on the partition
CREATE INDEX IF NOT EXISTS event_2025_11_entity_idx ON event_2025_11(entity_type, entity_id, ts DESC);
CREATE INDEX IF NOT EXISTS event_2025_11_type_ts_idx ON event_2025_11(event_type, ts DESC);

-- Note: Event table is partitioned, now we have a partition to insert into
INSERT INTO event (ts, entity_type, entity_id, event_type, data) VALUES
  (now() - interval '10 hours', 'job', '70000000-0000-0000-0000-000000000001', 'job.submitted', '{"owner_id":"00000000-0000-0000-0000-000000000001","priority":10}'::jsonb),
  (now() - interval '9 hours 55 minutes', 'job', '70000000-0000-0000-0000-000000000001', 'job.started', '{}'::jsonb),
  (now() - interval '9 hours 30 minutes', 'job', '70000000-0000-0000-0000-000000000001', 'job.completed', '{"duration_ms":1500000}'::jsonb),
  (now() - interval '5 hours', 'node', '50000000-0000-0000-0000-000000000001', 'node.session_started', '{"session_id":"60000000-0000-0000-0000-000000000001"}'::jsonb),
  (now() - interval '2 hours', 'job', '70000000-0000-0000-0000-000000000002', 'job.submitted', '{"owner_id":"00000000-0000-0000-0000-000000000002","priority":5}'::jsonb),
  (now() - interval '1 hour 55 minutes', 'job', '70000000-0000-0000-0000-000000000002', 'job.started', '{}'::jsonb),
  (now() - interval '30 minutes', 'job', '70000000-0000-0000-0000-000000000003', 'job.submitted', '{"owner_id":"00000000-0000-0000-0000-000000000003","priority":15}'::jsonb);

-- ========================================
-- 7) OUTBOX EVENTS (Some delivered, some pending)
-- ========================================

INSERT INTO outbox_event (created_at, topic, payload, delivered_at) VALUES
  (now() - interval '10 hours', 'metrics', '{"job_id":"70000000-0000-0000-0000-000000000001","event":"submitted"}'::jsonb, now() - interval '9 hours 59 minutes'),
  (now() - interval '9 hours 30 minutes', 'metrics', '{"job_id":"70000000-0000-0000-0000-000000000001","event":"completed"}'::jsonb, now() - interval '9 hours 29 minutes'),
  (now() - interval '2 hours', 'metrics', '{"job_id":"70000000-0000-0000-0000-000000000002","event":"submitted"}'::jsonb, now() - interval '1 hour 59 minutes'),
  (now() - interval '30 minutes', 'metrics', '{"job_id":"70000000-0000-0000-0000-000000000003","event":"submitted"}'::jsonb, NULL),
  (now() - interval '5 minutes', 'webhook', '{"node_id":"50000000-0000-0000-0000-000000000001","event":"task_completed","task_id":"80000000-0000-0000-0000-000000000012"}'::jsonb, NULL);

