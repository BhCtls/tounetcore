-- Add URL field to apps table
ALTER TABLE apps ADD COLUMN url TEXT DEFAULT '';

-- Add index on url field for better performance (optional)
CREATE INDEX IF NOT EXISTS idx_apps_url ON apps(url);

-- Update existing apps with empty URL (optional)
UPDATE apps SET url = '' WHERE url IS NULL;