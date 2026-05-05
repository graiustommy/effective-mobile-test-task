DROP TABLE IF EXISTS subscriptions CASCADE;
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(50),
    price INT,
    user_id UUID NOT NULL,     
	start_date VARCHAR(7) NOT NULL,
    end_date VARCHAR(7)
);

CREATE INDEX user_id_index ON subscriptions(user_id);
CREATE INDEX service_name_index ON subscriptions(service_name);
