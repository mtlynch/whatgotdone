import {mutations} from '@/store';

describe('mutations', () => {
  test('setUsername sets a username', () => {
    const state = {
      username: null,
    };
    mutations.setUsername(state, 'testUser123');
    expect(state.username).toBe('testUser123');
  });

  test('clearLoginState clears the local user state variables', () => {
    const state = {
      username: 'testUser123',
      following: ['alice', 'bob', 'charlie'],
      recentFollowingEntries: [{dummyEntry: true}, {dummyEntry: true}],
    };
    mutations.clearLoginState(state);
    expect(state.username).toBe(null);
    expect(state.following).toBe([]);
    expect(state.recentFollowingEntries).toBe([]);
  });

  test('setRecent adds recent entries', () => {
    const entries = [
      {
        key: '/testUser321/2019-09-27',
        author: 'testUser321',
        date: new Date(2019, 9, 27),
        markdown: 'I went to the zoo today',
      },
      {
        key: '/testUser456/2019-09-27',
        author: 'testUser456',
        date: new Date(2019, 9, 27),
        markdown: 'I ate an ice-cream sandwich yesterday',
      },
    ];
    const state = {
      recentEntries: null,
    };
    mutations.setRecent(state, entries);
    expect(state.recentEntries).toBe(entries);
  });
});
