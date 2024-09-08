import { Cell, List } from '@telegram-apps/telegram-ui';

import { Link } from '@/components/Link/Link.jsx';


import './IndexPage.css';

/**
 * @returns {JSX.Element}
 */
export function IndexPage() {
  const debug = false;
  return (
    <List>
      <Link to="/meetings_list/">
        <Cell subtitle="">Список совещаний</Cell>
      </Link>
      <Link to="/calendar/">
        <Cell subtitle="">Календарь</Cell>
      </Link>
      {debug &&<Link to="/test/">
        <Cell subtitle="test page">Api Test</Cell>
      </Link>}
    </List>
  );
}
