CREATE TABLE IF NOT EXISTS completions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID NOT NULL,
    routine_id UUID NOT NULL,
    occurrence_key VARCHAR(64) NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT unique_user_routine_occurrence UNIQUE (user_id, routine_id, occurrence_key)

);

CREATE INDEX IF NOT EXISTS idx_completions_user_id ON completions(user_id);

