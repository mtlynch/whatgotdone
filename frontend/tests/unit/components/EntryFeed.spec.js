import {shallowMount, createLocalVue} from '@vue/test-utils';

import BootstrapVue from 'bootstrap-vue';

import EntryFeed from '@/components/EntryFeed.vue';
import PartialJournal from '@/components/PartialJournal.vue';

const localVue = createLocalVue();
localVue.use(BootstrapVue);

describe('EntryFeed component', () => {
  test('renders entries as PartialJournals', () => {
    const mockEntries = [
      {
        permalink: '/testUser321/2019-09-27',
        author: 'testUser321',
        date: new Date(2019, 9, 27),
        markdown: 'I went to the zoo today',
      },
      {
        permalink: '/testUser456/2019-09-27',
        author: 'testUser456',
        date: new Date(2019, 9, 27),
        markdown: 'I ate an ice-cream sandwich yesterday',
      },
    ];

    const wrapper = shallowMount(EntryFeed, {
      propsData: {
        readEntriesFromStore: () => {
          return mockEntries;
        },
        readEntriesFromServer: () => {
          return new Promise(function (resolve) {
            resolve(mockEntries);
          });
        },
      },
      localVue,
    });

    const entries = wrapper.findAll(PartialJournal);
    expect(entries).toHaveLength(mockEntries.length);

    entries.wrappers.forEach((wrapper, i) => {
      expect(wrapper.vm.entry).toBe(mockEntries[i]);
    });
  });
});
