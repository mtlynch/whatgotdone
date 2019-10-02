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
