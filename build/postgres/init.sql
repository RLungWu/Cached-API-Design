CREATE TABLE "Ad" (
  "adId" bigserial PRIMARY KEY,
  "title" text,
  "startAt" timestamp,
  "endAt" timestamp,
  "conditionId" bigserial
);

CREATE TABLE "Condition" (
  "conditionId" bigserial,
  "ageStart" int,
  "country" text,
  "platform" text
);

CREATE TABLE "ConditionsCountry" (
  "conditionCountryId" bigserial PRIMARY KEY,
  "conditionId" bigserial,
  "country" text
);

CREATE TABLE "ConditionsPlatform" (
  "conditionPlatformId" bigserial PRIMARY KEY,
  "conditionId" bigserial,
  "platform" text
);

ALTER TABLE "Ad" ADD FOREIGN KEY ("conditionId") REFERENCES "Condition" ("conditionId");

ALTER TABLE "ConditionsCountry" ADD FOREIGN KEY ("conditionId") REFERENCES "Condition" ("conditionId");

ALTER TABLE "ConditionsPlatform" ADD FOREIGN KEY ("conditionId") REFERENCES "Condition" ("conditionId");
