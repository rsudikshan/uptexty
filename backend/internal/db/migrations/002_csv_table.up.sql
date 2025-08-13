CREATE TABLE csv_table(
id BIGSERIAL PRIMARY KEY,
file_name VARCHAR(200) NOT NULL,
uploaded_at TIMESTAMPTZ DEFAULT now(),
uploaded_by INT,
FOREIGN KEY (uploaded_by) REFERENCES user_table(id)
);



CREATE TABLE csv_rows (
  id BIGSERIAL PRIMARY KEY,
  csv_file_id BIGINT NOT NULL REFERENCES csv_table(id) ON DELETE CASCADE,
  position NUMERIC(30,10) NOT NULL,   -- high precision ordering key
  input_text TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- index to quickly fetch rows in order
CREATE INDEX idx_csv_rows_file_position ON csv_rows(csv_file_id, position);