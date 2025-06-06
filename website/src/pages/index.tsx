import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import CodeDemo from '@site/src/components/CodeDemo';
import QuickStart from '@site/src/components/QuickStart';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro">
            Get Started - 5min ⏱️
          </Link>
          <Link
            className={`button button--outline button--secondary button--lg ${styles.outlineButtonWhite}`}
            to="https://marketplace.visualstudio.com/items?itemName=keva-dev.go-ioc">
            VS Code Extension
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - Dependency Injection for Go`}
      description="Go IoC brings Spring-style autowiring to Go with compile-time dependency injection. Zero runtime overhead, type-safe, and familiar syntax for Java developers.">
      <HomepageHeader />
      <main>
        <CodeDemo />
        <HomepageFeatures />
        <QuickStart />
      </main>
    </Layout>
  );
}
