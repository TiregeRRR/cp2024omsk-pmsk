import { useLaunchParams } from '@telegram-apps/sdk-react';
import { List } from '@telegram-apps/telegram-ui';

import { DisplayData } from '@/components/DisplayData/DisplayData.jsx';

/**
 * @returns {JSX.Element}
 */
export function TestLaunchParamsPage() {
  const lp = useLaunchParams();

  return (
    <List>
      <DisplayData
        rows={[
          { title: 'tgWebAppPlatform', value: lp.platform },
          { title: 'tgWebAppShowSettings', value: lp.showSettings },
          { title: 'tgWebAppVersion', value: lp.version },
          { title: 'tgWebAppBotInline', value: lp.botInline },
          { title: 'tgWebAppStartParam', value: lp.startParam },
          { title: 'tgWebAppData', type: 'link', value: '/test/init-data' },
          { title: 'tgWebAppThemeParams', type: 'link', value: '/test/theme-params' },
        ]}
      />
    </List>
  );
}
