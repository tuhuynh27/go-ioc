import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import CodeBlock from '@theme/CodeBlock';
import styles from './styles.module.css';

const installCode = `go install github.com/tuhuynh27/go-ioc/cmd/iocgen@latest`;

export default function QuickStart(): ReactNode {
  return (
    <section className={styles.quickStart}>
      <div className="container">
        <div className="row">
          <div className="col col--8 col--offset-2">
            <div className="text--center margin-bottom--lg">
              <h2 style={{color: 'var(--ifm-heading-color)'}}>Get Started in Minutes</h2>
              <p style={{color: 'var(--ifm-color-content-secondary)'}}>Install Go IoC and start using dependency injection in your Go projects</p>
            </div>
            
            <div className={styles.installStep}>
              <h3>1. Install the CLI</h3>
              <CodeBlock language="bash">
                {installCode}
              </CodeBlock>
            </div>
            
            <div className={styles.steps}>
              <div className="row">
                <div className="col col--4">
                  <div className={styles.step}>
                    <div className={styles.stepNumber}>2</div>
                    <h4>Add Components</h4>
                    <p>Mark your structs with IoC annotations using struct tags</p>
                  </div>
                </div>
                <div className="col col--4">
                  <div className={styles.step}>
                    <div className={styles.stepNumber}>3</div>
                    <h4>Generate Code</h4>
                    <p>Run <code>iocgen</code> to generate type-safe wire code</p>
                  </div>
                </div>
                <div className="col col--4">
                  <div className={styles.step}>
                    <div className={styles.stepNumber}>4</div>
                    <h4>Use Container</h4>
                    <p>Initialize and use your dependency-injected services</p>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="text--center margin-top--lg">
              <Link
                className="button button--primary button--lg"
                to="/docs/intro">
                Read Full Documentation
              </Link>
              <Link
                className="button button--secondary button--lg margin-left--sm"
                to="https://github.com/tuhuynh27/go-ioc-gin-demo">
                View Example Project
              </Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}