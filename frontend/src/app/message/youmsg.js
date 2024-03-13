import styles from "../page.module.css";

export default function YouMsg() {
  return (
    <li className={styles.messageBox}>
      <div className={styles.username}>나</div>
      <div className={styles.message}></div>
    </li>
  );
}
