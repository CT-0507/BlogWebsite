function truncate(text: string, max = 100) {
  return text.length > max ? text.slice(0, max) + "..." : text;
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString();
}

export { truncate, formatDate };
