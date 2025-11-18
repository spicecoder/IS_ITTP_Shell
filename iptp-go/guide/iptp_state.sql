-- iptp_state.sql - Persistent Field State for gobash IPTP

-- Processes table: Named shell sessions
CREATE TABLE processes (
    process_name TEXT PRIMARY KEY,
    intention TEXT NOT NULL,
    current_dir TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT 1
);

-- Field table: Semantic pulses as first-class state
CREATE TABLE field (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    process_name TEXT NOT NULL,
    pulse_name TEXT NOT NULL,
    tv TEXT CHECK(tv IN ('Y', 'N', 'U')),
    response TEXT, -- JSON blob for structured data
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (process_name) REFERENCES processes(process_name)
);

-- Index for fast field matching
CREATE INDEX idx_field_lookup ON field(process_name, pulse_name, tv);

-- Intentions table: Communication between processes
CREATE TABLE intentions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    intention_name TEXT NOT NULL,
    source_process TEXT NOT NULL,
    target_process TEXT,
    signal TEXT, -- JSON array of pulses
    status TEXT CHECK(status IN ('pending', 'absorbed', 'reflected')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_process) REFERENCES processes(process_name)
);

-- Example queries:

-- Check if any process is deploying
SELECT * FROM field 
WHERE pulse_name = 'deployment_active' 
  AND tv = 'Y';

-- Get all pulses for a process (its current FIELD)
SELECT pulse_name, tv, response, timestamp
FROM field
WHERE process_name = 'auth'
ORDER BY timestamp DESC;

-- Find processes working on similar things
SELECT DISTINCT p1.process_name, p2.process_name
FROM field p1
JOIN field p2 ON p1.pulse_name = p2.pulse_name
WHERE p1.process_name != p2.process_name
  AND p1.tv = 'Y' AND p2.tv = 'Y';

-- Trigger conditions (field matching)
CREATE VIEW active_triggers AS
SELECT 
    i.intention_name,
    i.source_process,
    COUNT(*) as matched_pulses
FROM intentions i
JOIN field f ON f.process_name = i.target_process
WHERE json_extract(i.signal, '$[*].name') = f.pulse_name
  AND json_extract(i.signal, '$[*].TV') = f.tv
GROUP BY i.id;
