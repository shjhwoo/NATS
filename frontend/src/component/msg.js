
export default function Msg({username, content}) {
    return (
        <div class="message">
            <span class="username">{username}</span>
            <span class="content">{content}</span>
        </div> 
    )
}