const millisecondsPerWeek = 1000 * 60 * 60 * 24 * 7;

function millisecondsBetweenDates(a, b) {
  return a - b;
}

function millisecondsBetweenEntries(a, b) {
  return millisecondsBetweenDates(a.date, b.date);
}

function isEntryRecent(entry, currentTime) {
  const mostRecentEntry = entry;
  return (
    millisecondsBetweenDates(currentTime, mostRecentEntry.date) <=
    millisecondsPerWeek
  );
}

function areEntriesConsecutive(a, b) {
  return Math.abs(millisecondsBetweenEntries(a, b)) === millisecondsPerWeek;
}

export function latestStreak(entries, currentTime = new Date()) {
  if (entries.length == 0) {
    return 0;
  }

  const mostRecentEntry = entries[0];
  if (!isEntryRecent(mostRecentEntry, currentTime)) {
    return 0;
  }

  let streakLength = 1;
  for (let i = 1; i <= entries.length - 1; i++) {
    if (!areEntriesConsecutive(entries[i - 1], entries[i])) {
      return streakLength;
    }
    streakLength++;
  }
  return streakLength;
}

export function longestStreak(entries) {
  if (entries.length == 0) {
    return 0;
  }
  let longestStreak = 1;
  let currentStreak = 1;
  for (let i = 0; i <= entries.length - 2; i++) {
    if (areEntriesConsecutive(entries[i], entries[i + 1])) {
      currentStreak++;
    } else {
      currentStreak = 1;
    }
    longestStreak = Math.max(longestStreak, currentStreak);
  }
  return longestStreak;
}
