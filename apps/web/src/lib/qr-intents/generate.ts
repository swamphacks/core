export function generateCheckInIntent(userId: string, eventId: string) {
  return `checkin::${userId}+${eventId}`;
}
