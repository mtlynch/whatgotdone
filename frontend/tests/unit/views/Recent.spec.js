import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';

import Recent from '@/views/Recent.vue';
import PartialJournal from '@/components/PartialJournal.vue';

const localVue = createLocalVue();
localVue.use(Vuex);

describe('Recent.vue', () => {
  test('renders recent items as PartialJournals', () => {
    const recentEntries = [
      {
        key: `/testAuthor/01-01-2000`,
        author: {},
        date: new Date(2019),
        markdown: '',
      },
      {
        key: `/testAuthor/01-02-2000`,
        author: {},
        date: new Date(2019),
        markdown: '',
      },
    ];

    let state = {
      recentEntries,
    };

    const store = new Vuex.Store({
      state,
    });

    const wrapper = shallowMount(Recent, { store, localVue });

    const entries = wrapper.findAll(PartialJournal);
    expect(entries).toHaveLength(recentEntries.length);

    entries.wrappers.forEach((wrapper, i) => {
      expect(wrapper.vm.entry).toBe(recentEntries[i]);
    });
  });
});
