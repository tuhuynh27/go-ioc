import type {ReactNode} from 'react';
import styles from './styles.module.css';

interface TerminalProps {
  title?: string;
  children: ReactNode;
}

export default function Terminal({ title = "Terminal", children }: TerminalProps): ReactNode {
  return (
    <div className={styles.terminal}>
      <div className={styles.terminalHeader}>
        <div className={styles.trafficLights}>
          <div className={`${styles.light} ${styles.red}`}></div>
          <div className={`${styles.light} ${styles.yellow}`}></div>
          <div className={`${styles.light} ${styles.green}`}></div>
        </div>
        <div className={styles.terminalTitle}>{title}</div>
        <div className={styles.spacer}></div>
      </div>
      <div className={styles.terminalBody}>
        {children}
      </div>
    </div>
  );
}