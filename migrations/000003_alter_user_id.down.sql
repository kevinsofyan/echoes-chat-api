SET search_path TO echoes_chat;

-- Drop all foreign keys
ALTER TABLE tokens DROP CONSTRAINT IF EXISTS tokens_user_id_fkey;
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_created_by_fkey;
ALTER TABLE room_members DROP CONSTRAINT IF EXISTS room_members_user_id_fkey;
ALTER TABLE room_members DROP CONSTRAINT IF EXISTS room_members_room_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_sender_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_room_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_reply_to_id_fkey;

-- Convert back to INTEGER
ALTER TABLE users ALTER COLUMN id TYPE INTEGER USING 1;
ALTER TABLE users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

ALTER TABLE rooms ALTER COLUMN id TYPE INTEGER USING 1;
ALTER TABLE rooms ALTER COLUMN id SET DEFAULT nextval('rooms_id_seq'::regclass);
ALTER TABLE rooms ALTER COLUMN created_by TYPE INTEGER USING 1;

ALTER TABLE tokens ALTER COLUMN id TYPE INTEGER USING 1;
ALTER TABLE tokens ALTER COLUMN id SET DEFAULT nextval('tokens_id_seq'::regclass);
ALTER TABLE tokens ALTER COLUMN user_id TYPE INTEGER USING 1;

ALTER TABLE room_members ALTER COLUMN id TYPE INTEGER USING 1;
ALTER TABLE room_members ALTER COLUMN id SET DEFAULT nextval('room_members_id_seq'::regclass);
ALTER TABLE room_members ALTER COLUMN user_id TYPE INTEGER USING 1;
ALTER TABLE room_members ALTER COLUMN room_id TYPE INTEGER USING 1;

ALTER TABLE messages ALTER COLUMN id TYPE INTEGER USING 1;
ALTER TABLE messages ALTER COLUMN id SET DEFAULT nextval('messages_id_seq'::regclass);
ALTER TABLE messages ALTER COLUMN sender_id TYPE INTEGER USING 1;
ALTER TABLE messages ALTER COLUMN room_id TYPE INTEGER USING 1;
ALTER TABLE messages ALTER COLUMN reply_to_id TYPE INTEGER USING 1;

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