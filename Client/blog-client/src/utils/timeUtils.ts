export function relativeTime(createdAt: string) {
  const diff = (Date.now() - new Date(createdAt).getTime()) / 1000;

  if (diff < 60) return `${Math.floor(diff)}s`;

  const minutes = diff / 60;
  if (minutes < 60) return `${Math.floor(minutes)}m`;

  const hours = minutes / 60;
  if (hours < 24) return `${Math.floor(hours)}h`;

  const days = hours / 24;
  if (days < 30) return `${Math.floor(days)}d`;

  const months = days / 30;
  if (months < 12) return `${Math.floor(months)}mo`;

  const years = months / 12;
  return `${Math.floor(years)}y`;
}

/**
 * @example formatDayName("2026-06-01 00:00:00 +0000 UTC"); // returns "Monday"
 *
 * @param dateString
 * @returns
 */
export function formatDayName(dateString: string): string {
  const date = new Date(dateString);

  return new Intl.DateTimeFormat("en-US", {
    weekday: "long",
    timeZone: "UTC",
  }).format(date);
}

/**
 * @example formatDayShort("2026-06-01 00:00:00 +0000 UTC"); // returns "Mon"
 * @param dateString
 * @returns
 */
export function formatDayShort(dateString: string): string {
  const date = new Date(dateString);

  return new Intl.DateTimeFormat("en-US", {
    weekday: "short",
    timeZone: "UTC",
  }).format(date);
}
/**
 * @example formatWeekLabel("2026-06-01 00:00:00 +0000 UTC"); // returns "06/01"
 *
 * @param dateString
 * @returns
 */
export function formatWeekLabel(dateString: string): string {
  const date = new Date(dateString);

  return new Intl.DateTimeFormat("en-US", {
    month: "2-digit",
    day: "2-digit",
    timeZone: "UTC",
  }).format(date);
}
