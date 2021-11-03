-- creates tags table
CREATE TABLE tag (
    name text Primary Key,
    positive bool
);

-- sample insert queries in tags
INSERT INTO tag (name, positive) VALUES ('hardworking', true);
INSERT INTO tag (name, positive) VALUES ('slacker', false);
INSERT INTO tag (name, positive) VALUES ('leader', true);
INSERT INTO tag (name, positive) VALUES ('friendly', true);