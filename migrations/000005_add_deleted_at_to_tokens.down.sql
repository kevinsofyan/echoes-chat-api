SET search_path TO echoes_chat;

ALTER TABLE tokens DROP COLUMN IF EXISTS deleted_at;