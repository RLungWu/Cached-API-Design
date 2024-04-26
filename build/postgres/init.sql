CREATE TABLE "Ad" (
  "title" string PRIMARY KEY,
  "startAt" timestamp,
  "endAt" timestamp,
  "conditionId" bigserial
);

CREATE TABLE "Condition" (
  "conditionId" bigserial,
  "ageStart" int,
  "country" string,
  "platform" string
);

CREATE TABLE "ConditionsCountry" (
  "conditionCountryId" bigserial PRIMARY KEY,
  "conditionId" bigserial,
  "country" string
);

CREATE TABLE "ConditionsPlatform" (
  "conditionPlatformId" bigserial PRIMARY KEY,
  "conditionId" bigserial,
  "platform" string
);

ALTER TABLE "Ad" ADD FOREIGN KEY ("conditionId") REFERENCES "Condition" ("conditionId");

ALTER TABLE "ConditionsCountry" ADD FOREIGN KEY ("conditionId") REFERENCES "Condition" ("conditionId");

ALTER TABLE "ConditionsPlatform" ADD FOREIGN KEY ("conditionId") REFERENCES "Condition" ("conditionId");
