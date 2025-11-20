SET search_path TO echoes_chat;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop ALL foreign keys
ALTER TABLE tokens DROP CONSTRAINT IF EXISTS tokens_user_id_fkey;
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_created_by_fkey;
ALTER TABLE room_members DROP CONSTRAINT IF EXISTS room_members_user_id_fkey;
ALTER TABLE room_members DROP CONSTRAINT IF EXISTS room_members_room_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_sender_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_room_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_reply_to_id_fkey;

-- Convert users.id to UUID
ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
ALTER TABLE users ALTER COLUMN id TYPE UUID USING gen_random_uuid();
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Convert rooms.id to UUID
ALTER TABLE rooms ALTER COLUMN id DROP DEFAULT;
ALTER TABLE rooms ALTER COLUMN id TYPE UUID USING gen_random_uuid();
ALTER TABLE rooms ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Convert rooms.created_by to UUID
ALTER TABLE rooms ALTER COLUMN created_by TYPE UUID USING gen_random_uuid();

-- Convert tokens.id and user_id to UUID
ALTER TABLE tokens ALTER COLUMN id DROP DEFAULT;
ALTER TABLE tokens ALTER COLUMN id TYPE UUID USING gen_random_uuid();
ALTER TABLE tokens ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE tokens ALTER COLUMN user_id TYPE UUID USING gen_random_uuid();

-- Convert room_members IDs to UUID
ALTER TABLE room_members ALTER COLUMN id DROP DEFAULT;
ALTER TABLE room_members ALTER COLUMN id TYPE UUID USING gen_random_uuid();
ALTER TABLE room_members ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE room_members ALTER COLUMN user_id TYPE UUID USING gen_random_uuid();
ALTER TABLE room_members ALTER COLUMN room_id TYPE UUID USING gen_random_uuid();

-- Convert messages IDs to UUID
ALTER TABLE messages ALTER COLUMN id DROP DEFAULT;
ALTER TABLE messages ALTER COLUMN id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE messages ALTER COLUMN sender_id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages ALTER COLUMN room_id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages ALTER COLUMN reply_to_id TYPE UUID USING gen_random_uuid();

-- Add foreign keys back
ALTER TABLE tokens ADD CONSTRAINT tokens_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE rooms ADD CONSTRAINT rooms_created_by_fkey 
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE room_members ADD CONSTRAINT room_members_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE room_members ADD CONSTRAINT room_members_room_id_fkey 
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;

ALTER TABLE messages ADD CONSTRAINT messages_sender_id_fkey 
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE messages ADD CONSTRAINT messages_room_id_fkey 
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;

ALTER TABLE messages ADD CONSTRAINT messages_reply_to_id_fkey 
    FOREIGN KEY (reply_to_id) REFERENCES messages(id) ON DELETE SET NULL;