CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    group_id  varchar(50) NOT NULL,
    employee_id varchar(10) NOT NULL,
    organization_id varchar(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    pin varchar(50) NOT NULL,
    created_date TIMESTAMPTZ NOT NULL
);

CREATE TABLE access_tokens (
    token_id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    access_token TEXT NOT NULL,
    created_date TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE groups (
    id VARCHAR(50) PRIMARY KEY,
    group_key varchar(10) NOT NULL,
    name VARCHAR(50) NOT NULL,
    level SMALLINT
);

CREATE TABLE organizations (
    id VARCHAR(50) PRIMARY KEY,
    parent_id SMALLINT,
    group_id varchar(50)
);

-- ------------------
-- INITIALIZATION --
-- ------------------

-- Seeder untuk tabel groups
INSERT INTO groups (id, group_key, name, level) VALUES
    ('a7d39cf6-44b6-41fc-b3e9-7b16df5321c5', 'COM001', 'PT. Indodev Niaga Internet', 1),
    ('13bcb11c-111e-4a65-9afd-90a86a01cd21', 'BOD001', 'Board Of Directors', 2),
    ('201ddde1-f797-484b-b1a0-07d1190e790a', 'DVS001', 'Information Technology', 3),
    ('12f34634-00c2-45c9-b3fa-627f1b8634c6', 'TIM001', 'ERP Development', 4),
    ('bb6a88a1-c5b0-49b3-84b5-54ed4dbd9b84', 'EMP001', 'Employee', 5),
    ('216d379e-2be9-4ef2-9225-5face1fb0c5e', 'DVS002', 'Marketing and Sales', 3),
    ('b66e74e9-9f27-48e5-8387-37cd17ca7a20', 'TIM004', 'Sales', 4),
    ('6b5dcac0-2880-439a-b5eb-177a1e1c7013', 'TIM002', 'Tech Development', 4),
    ('6fd65892-0633-4cd2-917a-a17476bfafdb', 'TIM003', 'Software Maintenance', 4),
    ('16f98087-8d47-4984-b004-9fbf0f2b7e', 'TIM005', 'HR Development', 4);


INSERT INTO users (id, group_id, employee_id, organization_id, name, phone_number, pin, created_date) VALUES
     ('550e8400-e29b-41d4-a716-446655440001', '13bcb11c-111e-4a65-9afd-90a86a01cd21', 'EMP002', 'org2', 'Jane Doe', '081234567891', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440003', '12f34634-00c2-45c9-b3fa-627f1b8634c6', 'EMP004', 'org4', 'Bob', '081234567893', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440007', '201ddde1-f797-484b-b1a0-07d1190e790a', 'EMP008', 'org8', 'Frank', '081234567897', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440004', 'bb6a88a1-c5b0-49b3-84b5-54ed4dbd9b84', 'EMP005', '5', 'Charlie', '081234567894', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440009', 'bb6a88a1-c5b0-49b3-84b5-54ed4dbd9b84', 'EMP010', '5', 'Henry', '081234567899', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440006', 'b66e74e9-9f27-48e5-8387-37cd17ca7a20', 'EMP007', 'org7', 'Eva', '081234567896', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440002', '16f98087-8d47-4984-b004-9fbf0f2b7e', 'EMP003', 'org3', 'Alice', '081234567892', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440005', '6fd65892-0633-4cd2-917a-a17476bfafdb', 'EMP006', 'org6', 'David', '081234567895', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440000', 'b66e74e9-9f27-48e5-8387-37cd17ca7a20', 'EMP001', 'org1', 'John Tooer', '081234567890', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00'),
     ('550e8400-e29b-41d4-a716-446655440008', '6b5dcac0-2880-439a-b5eb-177a1e1c7013', 'EMP009', 'org9', 'Grace', '081234567898', '1a1dc91c907325c69271ddf0c944bc72', '2024-11-25 06:13:13.838212 +00:00');


INSERT INTO organizations (id, parent_id, group_id) VALUES
    ('1', null, 'a7d39cf6-44b6-41fc-b3e9-7b16df5321c5'),
    ('2', 1, '13bcb11c-111e-4a65-9afd-90a86a01cd21'),
    ('3', 2, '201ddde1-f797-484b-b1a0-07d1190e790a'),
    ('4', 3, '12f34634-00c2-45c9-b3fa-627f1b8634c6'),
    ('5', 4, 'bb6a88a1-c5b0-49b3-84b5-54ed4dbd9b84'),
    ('6', 2, '216d379e-2be9-4ef2-9225-5face1fb0c5e'),
    ('7', 6, 'b66e74e9-9f27-48e5-8387-37cd17ca7a20'),
    ('8', 3, '6b5dcac0-2880-439a-b5eb-177a1e1c7013'),
    ('9', 3, '6fd65892-0633-4cd2-917a-a17476bfafdb'),
    ('10', 3, '16f98087-8d47-4984-b004-9fbf0f2b7e');
