import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  icon: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Spring-like Syntax',
    icon: '🍃',
    description: (
      <>
        Familiar <code>@Autowired</code> and <code>@Component</code> syntax for Java developers. 
        Use struct tags and markers for clean, Spring-style dependency injection in Go.
      </>
    ),
  },
  {
    title: 'Compile-time Safety',
    icon: '⚡',
    description: (
      <>
        Zero runtime overhead with pure code generation. All dependencies are resolved 
        at compile time, ensuring type safety and optimal performance.
      </>
    ),
  },
  {
    title: 'Advanced Analysis',
    icon: '🔍',
    description: (
      <>
        Built-in dependency graph visualization, circular dependency detection, 
        and comprehensive component analysis tools for better architecture insights.
      </>
    ),
  },
];

function Feature({title, icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <div className={styles.featureIcon}>{icon}</div>
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

function AntiPatternNotice(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          <div className="col col--12">
            <div className="admonition admonition-caution">
              <div className="admonition-heading">
                <h5>⚠️ Important: Go Anti-Patterns & Migration Strategy</h5>
              </div>
              <div className="admonition-content">
                <p>
                  <strong>This library intentionally violates Go idioms</strong> to provide a familiar 
                  migration bridge for Java/Spring teams transitioning to Go. While technically functional, 
                  it introduces anti-patterns like magic struct tags and global state containers.
                </p>
                <p>
                  <strong>Use Case</strong>: Temporary productivity bridge for Java developers learning Go.
                  The compile-time approach ensures clean migration to idiomatic Go patterns when ready.
                </p>
                <p>
                  <a href="/docs/anti-patterns" className="button button--primary">
                    📖 Read Full Migration Strategy →
                  </a>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <>
      <AntiPatternNotice />
      <section className={styles.features}>
        <div className="container">
          <div className="row">
            {FeatureList.map((props, idx) => (
              <Feature key={idx} {...props} />
            ))}
          </div>
        </div>
      </section>
    </>
  );
}
