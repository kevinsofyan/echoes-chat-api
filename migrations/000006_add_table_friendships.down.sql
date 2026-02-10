SET search_path TO echoes_chat;

DROP TRIGGER IF EXISTS update_friendships_updated_at ON friendships;
DROP TABLE IF EXISTS friendships;