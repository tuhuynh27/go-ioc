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

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
