import { mutations } from '@/store';

describe('mutations', () => {
  test('setUsername sets a username', () => {
    const username = 'testUser123';
    const state = {
      username: null,
    };
    mutations.setUsername(state, username);
    expect(state.username).toBe(username);
  });

  test('clearUsername clears the username', () => {
    const state = {
      username: 'testUser123',
    };
    mutations.clearUsername(state);
    expect(state.username).toBe(null);
  });

  test('setRecent adds recent entries', () => {
    const entries = [
      {
        key: `/testAuthor/mm-dd-yyyy`,
        author: {},
        date: new Date(2019),
        markdown: '',
      },
      {
        key: `/testAuthor/mm-dd-yyyy`,
        author: {},
        date: new Date(2019),
        markdown: '',
      },
    ];
    const state = {
      recentEntries: null,
    };
    mutations.setRecent(state, entries);
    expect(state.recentEntries).toBe(entries);
  });
});
