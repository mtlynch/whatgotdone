import {latestStreak, longestStreak} from '@/controllers/Streaks.js';

describe('streaks controller', () => {
  test('latestStreak handles empty entries', () => {
    const entries = [];
    expect(latestStreak(entries)).toEqual(0);
  });

  test('latestStreak is one after a single, recent entry', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
    ];
    expect(
      latestStreak(entries, /* currentTime= */ new Date('2022-01-15'))
    ).toEqual(1);
  });

  test('latestStreak is zero after a single stale entry', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
    ];
    expect(
      latestStreak(entries, /* currentTime= */ new Date('2022-01-22'))
    ).toEqual(0);
  });

  test('latestStreak is two after two recent weeks of updates', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
      {
        date: new Date('2022-01-07'),
      },
    ];
    expect(
      latestStreak(entries, /* currentTime= */ new Date('2022-01-15'))
    ).toEqual(2);
  });

  test('latestStreak recognizes a break in consecutive updates', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
      {
        date: new Date('2022-01-07'),
      },
      {
        date: new Date('2021-12-24'),
      },
    ];
    expect(
      latestStreak(entries, /* currentTime= */ new Date('2022-01-15'))
    ).toEqual(2);
  });

  test('latestStreak includes the update for the current week', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
      {
        date: new Date('2022-01-07'),
      },
    ];
    expect(
      latestStreak(entries, /* currentTime= */ new Date('2022-01-12'))
    ).toEqual(2);
  });

  test('longestStreak handles empty entries', () => {
    const entries = [];
    expect(longestStreak(entries)).toEqual(0);
  });

  test('longestStreak is one after a single entry', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
    ];
    expect(longestStreak(entries)).toEqual(1);
  });

  test('longestStreak is two after two consecutive entries', () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
      {
        date: new Date('2022-01-07'),
      },
    ];
    expect(longestStreak(entries)).toEqual(2);
  });

  test("longestStreak finds the longest streak even when it's not the most recent", () => {
    const entries = [
      {
        date: new Date('2022-01-14'),
      },
      {
        date: new Date('2022-01-07'),
      },
      // Break in the streak on 2021-12-31
      {
        date: new Date('2021-12-24'),
      },
      {
        date: new Date('2021-12-17'),
      },
      {
        date: new Date('2021-12-10'),
      },
    ];
    expect(longestStreak(entries)).toEqual(3);
  });
});
