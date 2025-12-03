-- Creating a table of texts for reading
CREATE TABLE reading_texts (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    word_count INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Creating a table of questions for texts
CREATE TABLE reading_questions (
    id BIGSERIAL PRIMARY KEY,
    text_id BIGINT NOT NULL REFERENCES reading_texts(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    answer_option_1 TEXT NOT NULL,
    answer_option_2 TEXT NOT NULL,
    answer_option_3 TEXT NOT NULL,
    answer_option_4 TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for query optimization
CREATE INDEX idx_reading_texts_created_at ON reading_texts(created_at);
CREATE INDEX idx_reading_questions_text_id ON reading_questions(text_id);