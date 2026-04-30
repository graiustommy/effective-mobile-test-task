CREATE TABLE subscriptions IF NOT EXISTS (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid()
	service_name varchar(255) UNIQUE,
	price int,
	user_id UNIQUE UUID NOT NULL,
	start_date TIMESTAMP,
	end_date TIMESTAMP
);

CREATE INDEX user_id_service_name_index ON subscriptions(user_id, service_name);
