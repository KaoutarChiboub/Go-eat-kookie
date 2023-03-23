
---Creating a simple table about machines with corresponding id, name, outlet number and status

CREATE TABLE machine (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    outlet_number INT NOT NULL,
    status TEXT
);

INSERT INTO machine (name, outlet_number, status)
VALUES
  ('sopnode-l1', 3, 'Active'),
  ('n300', 11, 'Active'),
  ('sopnode-w3', 0, 'Inactive'),
  ('jaguar', 4, 'Active'),
  ('sopnode-z1', 9, 'Inactive');


