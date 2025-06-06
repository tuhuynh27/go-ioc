import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import styles from './styles.module.css';

interface ComparisonItem {
  feature: string;
  goIoc: string | ReactNode;
  wire: string | ReactNode;
  dig: string | ReactNode;
  inject: string | ReactNode;
}

const comparisonData: ComparisonItem[] = [
  {
    feature: 'Configuration Style',
    goIoc: 'Spring-like annotations',
    wire: 'Provider functions',
    dig: 'Constructor injection',
    inject: 'Field tags'
  },
  {
    feature: 'Runtime Overhead',
    goIoc: <span className={styles.excellent}>None</span>,
    wire: <span className={styles.excellent}>None</span>,
    dig: <span className={styles.poor}>Reflection-based</span>,
    inject: <span className={styles.poor}>Reflection-based</span>
  },
  {
    feature: 'Compile-time Safety',
    goIoc: <span className={styles.excellent}>✅ Full</span>,
    wire: <span className={styles.excellent}>✅ Full</span>,
    dig: <span className={styles.warning}>⚠️ Partial</span>,
    inject: <span className={styles.poor}>❌ None</span>
  },
  {
    feature: 'Auto Component Scanning',
    goIoc: <span className={styles.excellent}>✅ Yes</span>,
    wire: <span className={styles.poor}>❌ Manual</span>,
    dig: <span className={styles.poor}>❌ Manual</span>,
    inject: <span className={styles.poor}>❌ Manual</span>
  },
  {
    feature: 'Qualifier Support',
    goIoc: <span className={styles.excellent}>✅ Built-in</span>,
    wire: <span className={styles.poor}>❌ None</span>,
    dig: <span className={styles.warning}>⚠️ Limited</span>,
    inject: <span className={styles.poor}>❌ None</span>
  },
  {
    feature: 'Dependency Analysis',
    goIoc: <span className={styles.excellent}>✅ Advanced</span>,
    wire: <span className={styles.poor}>❌ None</span>,
    dig: <span className={styles.poor}>❌ Basic</span>,
    inject: <span className={styles.poor}>❌ None</span>
  },
  {
    feature: 'Learning Curve',
    goIoc: <span className={styles.excellent}>Low (Spring devs)</span>,
    wire: 'Medium',
    dig: 'Medium',
    inject: 'Low'
  }
];

export default function ComparisonTable(): ReactNode {
  return (
    <section className={styles.comparison}>
      <div className="container">
        <div className="row">
          <div className="col col--12">
            <div className="text--center margin-bottom--lg">
              <h2>Why Choose Go IoC?</h2>
              <p className="hero__subtitle">
                See how Go IoC compares to other dependency injection libraries
              </p>
            </div>
            
            <div className={styles.tableWrapper}>
              <table className={styles.comparisonTable}>
                <thead>
                  <tr>
                    <th>Feature</th>
                    <th className={styles.highlight}>Go IoC</th>
                    <th>Google Wire</th>
                    <th>Uber Dig</th>
                    <th>Facebook Inject</th>
                  </tr>
                </thead>
                <tbody>
                  {comparisonData.map((item, index) => (
                    <tr key={index}>
                      <td className={styles.featureCell}>{item.feature}</td>
                      <td className={`${styles.goIocCell} ${styles.highlight}`}>{item.goIoc}</td>
                      <td>{item.wire}</td>
                      <td>{item.dig}</td>
                      <td>{item.inject}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            
            <div className={styles.highlights}>
              <div className="row">
                <div className="col col--4">
                  <div className={styles.highlightBox}>
                    <h4>🍃 Spring-like Syntax</h4>
                    <p>Familiar <code>@Component</code> and <code>@Autowired</code> experience for Java developers</p>
                  </div>
                </div>
                <div className="col col--4">
                  <div className={styles.highlightBox}>
                    <h4>⚡ Zero Runtime Cost</h4>
                    <p>Pure compile-time code generation with no reflection overhead</p>
                  </div>
                </div>
                <div className="col col--4">
                  <div className={styles.highlightBox}>
                    <h4>🔍 Advanced Analysis</h4>
                    <p>Built-in dependency graph visualization and circular dependency detection</p>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="text--center margin-top--lg">
              <Link
                className="button button--primary button--lg"
                to="/docs/comparison">
                View Detailed Comparison
              </Link>
              <Link
                className="button button--secondary button--lg margin-left--sm"
                to="/docs/intro">
                Get Started Now
              </Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}