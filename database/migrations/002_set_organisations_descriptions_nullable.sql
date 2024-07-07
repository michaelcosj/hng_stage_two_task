-- Write your migrate up statements here
ALTER TABLE organisations ALTER COLUMN description DROP NOT NULL;

---- create above / drop below ----
Update organisations Set description="" where description is NULL;
ALTER TABLE organisations ALTER COLUMN description TEXT NOT NULL;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
