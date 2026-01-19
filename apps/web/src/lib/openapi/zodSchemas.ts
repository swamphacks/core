/**
 * Zod schemas for OpenAPI specification components.
 * This file needs to be manually kept in sync with the OpenAPI schema definitions.
 */

import { z } from "zod";

export const PlatformRoleSchema = z.enum(["user", "superuser"]);
export const EventRoleSchema = z.enum([
  "attendee",
  "admin",
  "staff",
  "applicant",
]);
