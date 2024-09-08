import { CalendarPage } from '@/pages/CalendarPage/CalendarPage';
import { IndexPage } from '@/pages/IndexPage/IndexPage';
import { MeetingsPage } from '@/pages/MeetingsPage/MeetingsPage';
import { MeetingsListPage } from '@/pages/MeetingsListPage/MeetingsListPage';
import { TestIndexPage as TestIndexPage } from '@/pages/TestIndexPage/TestIndexPage';
import { TestInitDataPage } from '@/pages/TestInitDataPage/TestInitDataPage';
import { TestLaunchParamsPage } from '@/pages/TestLaunchParamsPage/TestLaunchParamsPage.jsx';
import { TestThemeParamsPage } from '@/pages/TestThemeParamsPage/TestThemeParamsPage.jsx';

/**
 * @typedef {object} Route
 * @property {string} path
 * @property {import('react').ComponentType} Component
 * @property {string} [title]
 * @property {import('react').JSX.Element} [icon]
 */

/**
 * @type {Route[]}
 */
export const routes = [
  { path: '/', Component: IndexPage},
  { path: '/meetings/:id', Component: MeetingsPage},
  { path: '/calendar', Component: CalendarPage},
  { path: '/meetings_list', Component: MeetingsListPage},


  { path: '/test', Component: TestIndexPage },
  { path: '/test/init-data', Component: TestInitDataPage, title: 'Init Data' },
  { path: '/test/theme-params', Component: TestThemeParamsPage, title: 'Theme Params' },
  { path: '/test/launch-params', Component: TestLaunchParamsPage, title: 'Launch Params' },
];
