CREATE TABLE VEHICLE (
    vehicleid VARCHAR(6),
    licensenumber VARCHAR(10),
    statecode VARCHAR(2),
    PRIMARY KEY(vehicleid)
);

CREATE TABLE PERMIT (
    ownername VARCHAR(30),
    permitid VARCHAR(8),
    PRIMARY KEY(permitid)
);

CREATE TABLE VEH_PER (
    vid VARCHAR(6),
    pid VARCHAR(8),
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_p(pid) REFERENCES PERMIT(permitid),
    PRIMARY KEY(vid,pid)
);