import React from "react";
import { Cell, List } from '@telegram-apps/telegram-ui';
import { DateCalendar } from '@mui/x-date-pickers/DateCalendar';
import { Link } from '@/components/Link/Link.jsx';
import dayjs from 'dayjs';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { PickersDay } from '@mui/x-date-pickers/PickersDay';
import Badge from '@mui/material/Badge';
import Menu from '@mui/material/Menu';
import './CalendarPage.css';
import { useThemeParams } from '@telegram-apps/sdk-react';
import { createMuiTheme } from "../../functions/createMuiTheme";
import { ThemeProvider } from '@mui/material';
import { useState, useEffect } from "react";
import {host} from '@/host';
 
function ListInfo({ data }) {
  return (
    <>
      <List>
        {
          data.map((d) =>
            // ToDo ссылки на конкретное совещание
            <Link key={d.id} to={"/meetings/" + d.id}>
              <Cell>
                {d.created_at.toLocaleString()} {d.name}
              </Cell>
            </Link>
          )
        }
      </List>
    </>
  )
}

function Day(props) {
  const { meets = [], day, outsideCurrentMonth, ...other } = props;

  const daysData = meets.filter((inf) => dayjs(inf.created_at).month() === props.day.month());

  const highlightedDays = daysData.map((d) => dayjs(d.created_at).date())
  const currentData = daysData.filter((d) => dayjs(d.created_at).date() === props.day.date());

  const isSelected =
    !props.outsideCurrentMonth && highlightedDays.indexOf(props.day.date()) >= 0;

  // menu
  const [anchorEl, setAnchorEl] = React.useState(null);

  const open = Boolean(anchorEl);
  const handleClick = (event) => {
    if (isSelected) {
      setAnchorEl(event.currentTarget);
    }
  }

  const handleClose = () => {
    setAnchorEl(null);
  }

  return (
    <Badge
      key={props.day.toString()}
      overlap="circular"
      sx={{ cursor: "pointer" }}
      badgeContent={isSelected ? '📃' : undefined}
    >
      <PickersDay onClick={handleClick} {...other} outsideCurrentMonth={outsideCurrentMonth} day={day} />

      <Menu
        id="basic-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        MenuListProps={{
          'aria-labelledby': 'basic-button',
        }}
      >
        <ListInfo
          data={currentData}
        />
      </Menu>
    </Badge>
  );
}

/**
 * @returns {JSX.Element}
 */
export function CalendarPage() {
  //const listInfo = getListInfoMeeting();
  const tgTheme = useThemeParams();

  const [error, setError] = useState(undefined)
  const [meets, setMeets] = useState([])

  const theme = createMuiTheme(tgTheme.getState());

  useEffect(() => {
    fetch(host + "/get_transcriptions", { headers: { "ngrok-skip-browser-warning": '1' } })
      .then(response => {
        setError(undefined)
        if (response.ok && response.status == 200) {
          return response.json();
        }
        throw new Error('Не удалось совершить запрос на сервер')
      })
      .then(json => {
        setMeets(json.map((j) => { return { ...j, created_at: new Date(j.created_at) } }))
      })
      .catch(e => setError(e.message))
  }, [])

  if (error) {
    return <div>Ошибка - {error}</div>
  }
  return (
    <ThemeProvider theme={theme}>
      <LocalizationProvider dateAdapter={AdapterDayjs}>
        <DateCalendar
          defaultValue={dayjs()}
          slots={{
            day: Day,
          }}
          slotProps={{
            day: {
              meets,
            },
          }}
        />
      </LocalizationProvider>
    </ThemeProvider>
  );
}
