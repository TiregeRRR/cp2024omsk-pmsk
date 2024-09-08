import { Section, Cell, List, Button, Input } from '@telegram-apps/telegram-ui';

import { Link } from '@/components/Link/Link.jsx';


import './TestIndexPage.css';

import { useEffect, useState } from 'react';

/**
 * @returns {JSX.Element}
 */
export function TestIndexPage() {
  const [ds, setDs] = useState(96)

  return (
    <List>
      <Section
        header="Транскрибатершная"
        footer="Тестовая версия, что-то может быть явно не тем, чем является"
      >
        <Button>123321</Button>
        <Input type='number' value={ds} placeholder='dick size' onChange={e => setDs(e.target.value)}/>
        
        <Link to="/test/init-data">
          <Cell subtitle="User data, chat information, technical data">Init Data</Cell>
        </Link>
        <Link to="/test/launch-params">
          <Cell subtitle="Platform identifier, Mini Apps version, etc.">Launch Parameters</Cell>
        </Link>
        <Link to="/test/theme-params">
          <Cell subtitle="Telegram application palette information">Theme Parameters </Cell>
        </Link>
      </Section>
    </List>
  );
}
