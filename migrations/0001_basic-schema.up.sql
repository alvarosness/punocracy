ALTER DATABASE CHARACTER SET "utf8";

DROP TABLE IF EXISTS Permissions_T;
CREATE TABLE Permissions_T(
    permLevel INT NOT NULL,
    permDescription VARCHAR(30),

    CONSTRAINT Permissions_PK PRIMARY KEY (permLevel)
);

DROP TABLE IF EXISTS Users_T;
CREATE TABLE Users_T(
    userID INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    passwordHash VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    permLevel INT NOT NULL,

    CONSTRAINT Users_PK PRIMARY KEY (userID),
    CONSTRAINT Users_FK FOREIGN KEY (permLevel) REFERENCES Permissions_T(permLevel)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);

DROP TABLE IF EXISTS Words_T;
CREATE TABLE Words_T(
    wordID INT NOT NULL,
    word VARCHAR(20) NOT NULL UNIQUE,
    homophoneGroup INT NOT NULL,

    CONSTRAINT Words_PK PRIMARY KEY (wordID)
);
