-- Drop the foreign key constraint if it exists
ALTER TABLE IF EXISTS "account" DROP CONSTRAINT IF EXISTS "account_owner_fkey";

-- Drop the "users" table if it exists
DROP TABLE IF EXISTS "users";
