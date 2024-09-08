import { Section, List, Timeline, Button, Headline, Spinner, Title, Text, Snackbar } from '@telegram-apps/telegram-ui';

import { useState, useEffect } from 'react';
import { Link } from '@/components/Link/Link.jsx';
import WavesurferPlayer from '@wavesurfer/react'

import './MeetingsPage.css';
import { useParams } from 'react-router-dom';
import { host } from '@/host';


const formatTime = function (time) {
  if (!time) {
    return ''
  }
  return [
    Math.floor((time % 3600) / 60), // minutes
    ('00' + Math.floor(time % 60)).slice(-2) // seconds
  ].join(':')
}



function Player({ src }) {
  const [wavesurfer, setWavesurfer] = useState(null)
  const [isPlaying, setIsPlaying] = useState(false)
  const [isReady, setIsReady] = useState(false)

  const [backedsrc, setbackersrc] = useState(null)

  const [time, setTime] = useState(0)

  const onReady = (ws) => {
    setWavesurfer(ws)
    setIsPlaying(false)
    setIsReady(true)
  }

  const onPlayPause = () => {
    wavesurfer && wavesurfer.playPause()
  }

  const btn_prefix = isPlaying ? 'Остановить' : 'Воспроизвести'

  const time_str = time ? `${formatTime(time)}/${formatTime(wavesurfer.getDuration())}` : ''

  useEffect(() => {
    fetch(host + src, { headers: { "ngrok-skip-browser-warning": '1' } })
      .then(d => d.blob())
      .then(data => {
        const newUrl = URL.createObjectURL(data)
        console.debug(newUrl)
        setbackersrc(newUrl)
      })


  }, [])

  return (
    <div>
      <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center' }}>

        <Headline>Запись совещания:</Headline>


        {isReady && <Button onClick={onPlayPause}>
          {btn_prefix} {time_str}
        </Button>}

        {!isReady && <Spinner size='l' />}
      </div>


      {backedsrc && <WavesurferPlayer

        height={75}
        waveColor="#6ab3f3"
        url={backedsrc}

        autoScroll={true}
        autoCenter={true}
        hideScrollbar={true}
        minPxPerSec={25}
        onReady={onReady}
        onTimeupdate={(w) => setTime(w.getCurrentTime())}
        onPlay={() => setIsPlaying(true)}
        onPause={() => setIsPlaying(false)}
      />}
    </div>)


}

/**
 * @returns {JSX.Element}
 */
export function MeetingsPage() {

  const { id } = useParams()
  const [showMessage, setShowMessage] = useState(false)
  const [error, setError] = useState(undefined)
  const [meet, setMeet] = useState(undefined)



  const download = (format, type) => {
    const q = host + "/send_report/" + format + "/" + type + "/" + id
    fetch(q, { headers: { "ngrok-skip-browser-warning": '1' } })
    setShowMessage(true)
  }



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
        const items = json.filter(j => j.id == id).map((j) => { return { ...j, created_at: new Date(j.created_at) } })
        if (items.length == 0) {
          setError("Встреча не найдена")
        }
        setMeet(items[0])
      })
      .catch(e => setError(e.message))
  }, [])


  if (error) {
    return <div>Ошибка - {error}</div>
  }
  if (!meet) {
    return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', padding: 10, height: '30vh' }}>
      <Spinner />
    </div>
  }
  return (
    <List>
      <Section
        header={<div style={{ display: 'flex', justifyContent: 'space-around', padding: 10 }}>
          <Text>{meet.name}</Text>
          <Text>{meet.created_at.toUTCString()}</Text>
        </div>}
        footer=""
      >
        {meet.status < 4 && <Timeline
          active={meet.status}
        >
          <Timeline.Item header="Загрузка">

          </Timeline.Item>
          <Timeline.Item header="Транскрибация">
            Определение речи
          </Timeline.Item>
          <Timeline.Item header="Диаризация">
            Распределение текста по участникам
          </Timeline.Item>
          <Timeline.Item header="Генерация отчетов">
          </Timeline.Item>
          <Timeline.Item header="Сохранение">

          </Timeline.Item>

        </Timeline>}

        <Player src={meet.audio_link} />
        <p></p>

        {meet.status == 4 && <>
          <div className='flexrow'>
            <Headline>Официальный отчет</Headline>
          </div>
          <div className='flexrow'>
            <Button onClick={() => download("docx", "official")}>
              Скачать .docx
            </Button>
            <Button onClick={() => download("pdf", "official")}>
              Скачать .pdf
            </Button>
          </div>

          <div className='flexrow'>
            <Headline>Нефициальный отчет</Headline>
          </div>
          <div className='flexrow'>
            <Button onClick={() => download("docx", "unofficial")}>
              Скачать .docx
            </Button>
            <Button onClick={() => download("pdf", "unofficial")}>
              Скачать .pdf
            </Button>
          </div>
        </>
        }

        {
          meet.status == 5 && <div className='flexrow'><Headline>Ошибка при обработке совещания</Headline></div>
        }

      </Section>

      {showMessage && <Snackbar onClose={() => setShowMessage(false)} children={"Загрузка файла"} duration={5 * 1000}>
        Файл был отправлен личным сообщением через бота, сверните приложение
      </Snackbar>}


    </List>
  );
}
