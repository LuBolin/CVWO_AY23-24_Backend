CREATE TABLE Users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE Passwords (
    user_id INT NOT NULL,
    hashed VARBINARY(255) NOT NULL,
    salt VARBINARY(255) NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

CREATE TABLE Posts (
    id INT NOT NULL AUTO_INCREMENT,
    author_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (author_id) REFERENCES Users(id)
);

CREATE TABLE Comments (
    id INT NOT NULL AUTO_INCREMENT,
    author_id INT NOT NULL,
    post_id INT NOT NULL,
    content TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (author_id) REFERENCES Users(id),
    FOREIGN KEY (post_id) REFERENCES Posts(id)
);