import "./MeetingsListPage.css";
import { List, Cell, Text, Spinner } from '@telegram-apps/telegram-ui';
import { Link } from '@/components/Link/Link.jsx';
import { useEffect, useState } from "react";
import {host} from '@/host';
 
export function MeetingsListPage() {
    const [error, setError] = useState(undefined)
    const [meets, setMeets] = useState([])

    
    useEffect(() => {
         fetch(host + "/get_transcriptions", {headers: {"ngrok-skip-browser-warning": '1'}})
            .then(response => {
                setError(undefined)
                if (response.ok && response.status == 200) {
                    return response.json();
                }
                throw new Error('Не удалось совершить запрос на сервер')
            })
            .then(json => {
                setMeets(json.map((j) => {return {...j, created_at: new Date(j.created_at)}}))
            })
            .catch(e => setError(e.message) )
    }, [])


    if (error) {
        return <div>Ошибка - {error}</div>
    }
    return (
        // ToDo: Можно добавить фильтрацию и сортировку (по дате / названию)
        <List>
            {
                meets.map((meet) =>
                    // ToDo подставить в Link ссылку на страницу с конкретным совещанием
                    <Link key={meet.id} to={"/meetings/" + meet.id} className="container_item_meet" >
                        <Cell className="container_item_meet_info">
                            <Text weight="3">
                                {meet.name}
                            </Text>
                            <br />
                            <Text weight="3" >
                                {meet.created_at.toUTCString()}   {meet.status < 4 && <>- готовится</>} {meet.status == 4 && <>- готово</>} {meet.status == 5 && <>- ошибка</>}
                            </Text>
                             
                                   
                            

                        </Cell>
                    </Link>
                )
            }
        </List>
    )
}