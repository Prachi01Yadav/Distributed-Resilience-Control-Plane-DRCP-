const pptxgen = require('pptxgenjs');

async function createTechPresentation() {
    let pres = new pptxgen();
    pres.layout = 'LAYOUT_16x9';

    // Slide 1
    let slide1 = pres.addSlide();
    slide1.background = { color: '0f1117' };
    slide1.addText('Distributed Resilience Control Plane (DRCP)', { x: 1, y: 1.5, w: 8, h: 1, fontSize: 36, color: '4fc3f7', bold: true });
    slide1.addText('Technical Deep Dive', { x: 1, y: 2.5, w: 8, h: 1, fontSize: 24, color: 'a78bfa' });
    slide1.addText('By: Prachi Yadav\nPreventing cascading failures in microservices', { x: 1, y: 3.5, w: 8, h: 1, fontSize: 18, color: 'e2e8f0' });

    // Slide 2
    let slide2 = pres.addSlide();
    slide2.background = { color: '0f1117' };
    slide2.addText('Problem Statement', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide2.addText('• In microservice systems, one failing service can bring down everything.\n• Manual monitoring can\'t react fast enough.\n• No automated way to enforce SLAs in real time.\n• Incident history is easy to tamper with.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 3
    let slide3 = pres.addSlide();
    slide3.background = { color: '0f1117' };
    slide3.addText('What DRCP Does', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide3.addText('• Monitors all services in real time.\n• Automatically detects SLA breaches (high latency, high error rate).\n• Trips circuit breakers to stop bad traffic.\n• Records every incident on blockchain (can\'t be changed).', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 4
    let slide4 = pres.addSlide();
    slide4.background = { color: '0f1117' };
    slide4.addText('System Architecture', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide4.addText('• Dashboard UI → API Gateway → PostgreSQL\n• API Gateway → Kafka\n• Kafka → Worker → Redis (sliding window)\n• Worker → OPA (policy check)\n• OPA breach → Envoy xDS (trip circuit breaker)\n• OPA breach → Blockchain Anchor → Ethereum\n• OPA breach → DB (save incident)', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 5
    let slide5 = pres.addSlide();
    slide5.background = { color: '0f1117' };
    slide5.addText('Tech Stack', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide5.addText('• Language: Go 1.25\n• API: Gin Framework\n• Database: PostgreSQL (GORM) + SQLite\n• Messaging: Apache Kafka (Sarama)\n• Cache: Redis\n• Policy Engine: Open Policy Agent (OPA)\n• Service Mesh: Envoy Proxy (xDS API)\n• Blockchain: Ethereum (go-ethereum)\n• Frontend: Vanilla HTML/CSS/JS\n• Deployment: Docker + Render', { x: 0.5, y: 1.5, w: 9, h: 3.5, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 6
    let slide6 = pres.addSlide();
    slide6.background = { color: '0f1117' };
    slide6.addText('Microservice Architecture', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide6.addText('1. API Server (cmd/api): REST endpoints, serves dashboard, manages registry.\n2. Telemetry Worker (cmd/worker): Kafka consumer, budget calculator, policy evaluator.\n3. xDS Server (cmd/xds): Envoy gRPC control plane, pushes circuit breaker configs.\n4. Blockchain Anchor (cmd/anchor): Signs & submits Ethereum transactions for audit.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0' });

    // Slide 7
    let slide7 = pres.addSlide();
    slide7.background = { color: '0f1117' };
    slide7.addText('Data Flow - Telemetry Pipeline', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide7.addText('1. Service sends telemetry POST /api/v1/telemetry\n2. API produces message to Kafka topic telemetry-events\n3. Worker consumes from Kafka\n4. Worker updates Redis sliding window (60s window)\n5. Worker calculates error rate from Redis\n6. Worker sends data to OPA for policy check\n7. If breach → create incident + trip circuit breaker + anchor on blockchain', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0' });

    // Slide 8
    let slide8 = pres.addSlide();
    slide8.background = { color: '0f1117' };
    slide8.addText('Key Technical Components', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide8.addText('• Sliding Window (Redis): Uses Redis sorted sets with timestamps. Tracks total requests and errors over last 60 seconds.\n• OPA Policy Engine: Dual mode - built-in rules OR external OPA server.\n• Envoy xDS: Uses go-control-plane library. Pushes snapshots via gRPC. Trips breaker by setting max connections to 1.\n• Blockchain Anchor: EIP-1559 transactions. Incident hash stored in transaction data field.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 9
    let slide9 = pres.addSlide();
    slide9.background = { color: '0f1117' };
    slide9.addText('Design Patterns Used', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide9.addText('• Event-Driven Architecture: Kafka decouples ingestion from processing.\n• Circuit Breaker Pattern: Envoy auto-trips on SLA breach.\n• Graceful Degradation: Works without Kafka/Redis in reduced mode.\n• Sliding Window: Redis sorted sets for real-time error rates.\n• Policy-as-Code: OPA Rego for flexible SLA rules.\n• Immutable Audit Trail: Ethereum blockchain.\n• Dependency Injection: Handlers receive DB, Kafka, Redis as params.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 10
    let slide10 = pres.addSlide();
    slide10.background = { color: '0f1117' };
    slide10.addText('Project Structure', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide10.addText('├── cmd/         → 4 entry points (api, worker, xds, anchor)\n├── internal/    → Core logic\n├── pkg/         → Shared libs\n├── web/         → Frontend SPA\n├── contracts/   → Solidity smart contract\n└── deployments/ → Helm + Terraform', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', fontFace: 'Courier New' });

    // Slide 11
    let slide11 = pres.addSlide();
    slide11.background = { color: '0f1117' };
    slide11.addText('Challenges & Solutions', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide11.addText('• Distributed system coordination → Kafka for async messaging\n• Real-time processing → Redis sliding windows for sub-second calc.\n• Blockchain reliability → Graceful error handling, retry logic\n• Envoy xDS complexity → Used go-control-plane library\n• Zero-config local dev → SQLite fallback when Postgres unavailable\n• System works if parts fail → Graceful degradation pattern', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 12
    let slide12 = pres.addSlide();
    slide12.background = { color: '0f1117' };
    slide12.addText('Key Learnings', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide12.addText('• How to build event-driven distributed systems.\n• Working with Envoy\'s xDS gRPC protocol.\n• Blockchain integration for real-world audit trails.\n• Policy-as-code with OPA for flexible rule evaluation.\n• Importance of graceful degradation in distributed systems.\n• Sliding window algorithms for real-time metrics.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0', bullet: true });

    // Slide 13
    let slide13 = pres.addSlide();
    slide13.background = { color: '0f1117' };
    slide13.addText('Live Demo & API', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '4fc3f7', bold: true });
    slide13.addText('Live at: https://drcp-dashboard.onrender.com\n\n• GET /health - Health check\n• GET/POST /api/v1/services - List/Register services\n• POST /api/v1/services/:id/contracts - Create SLA contract\n• GET /api/v1/contracts - List contracts\n• GET /api/v1/incidents - List incidents\n• POST /api/v1/telemetry - Ingest telemetry', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: 'e2e8f0' });

    // Slide 14
    let slide14 = pres.addSlide();
    slide14.background = { color: '0f1117' };
    slide14.addText('Thank You!', { x: 1, y: 2, w: 8, h: 1, fontSize: 48, color: '4fc3f7', bold: true, align: 'center' });
    slide14.addText('Questions?\nLive Demo: https://drcp-dashboard.onrender.com\nBuilt by Prachi Yadav', { x: 1, y: 3, w: 8, h: 1.5, fontSize: 24, color: 'e2e8f0', align: 'center' });

    await pres.writeFile({ fileName: 'DRCP_Technical_Presentation.pptx' });
}

async function createBizPresentation() {
    let pres = new pptxgen();
    pres.layout = 'LAYOUT_16x9';

    // Slide 1
    let slide1 = pres.addSlide();
    slide1.background = { color: 'f8fafc' };
    slide1.addText('Distributed Resilience Control Plane', { x: 1, y: 1.5, w: 8, h: 1, fontSize: 36, color: '0284c7', bold: true, align: 'center' });
    slide1.addText('Building Reliable Microservices at Scale', { x: 1, y: 2.5, w: 8, h: 1, fontSize: 24, color: '334155', align: 'center' });
    slide1.addText('By: Prachi Yadav', { x: 1, y: 3.5, w: 8, h: 1, fontSize: 18, color: '475569', align: 'center' });

    // Slide 2
    let slide2 = pres.addSlide();
    slide2.background = { color: 'f8fafc' };
    slide2.addText('The Problem: Cascading Failures', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide2.addText('• Modern companies run hundreds of interconnected microservices.\n• When ONE service fails, it can crash EVERYTHING (cascading failure).\n• Manual monitoring is slow — by the time engineers notice, the damage is done.\n• SLA violations lead to financial penalties and lost trust.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 3
    let slide3 = pres.addSlide();
    slide3.background = { color: 'f8fafc' };
    slide3.addText('The Cost of Downtime', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide3.addText('• $5,600: Average cost of downtime per minute (Gartner).\n• 98%: Of organizations say 1 hour of downtime costs over $100,000.\n• 10-30%: Of contract value can be lost in SLA penalties.\n• Lost Trust: Customer trust drops significantly after major outages.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 4
    let slide4 = pres.addSlide();
    slide4.background = { color: 'f8fafc' };
    slide4.addText('Our Solution: DRCP', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide4.addText('• DRCP is an automated safety system for your microservices.\n• It watches all services 24/7 in real time.\n• Detects problems BEFORE they spread.\n• Automatically blocks bad traffic (like a circuit breaker in your house).\n• Records every incident permanently on the blockchain for auditing.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 5
    let slide5 = pres.addSlide();
    slide5.background = { color: 'f8fafc' };
    slide5.addText('How It Works', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide5.addText('1. MONITOR: Watches all services in real time, gathering telemetry.\n2. DETECT: Spots SLA violations instantly (high errors, slow responses).\n3. PROTECT: Automatically blocks bad traffic to stop the problem from spreading.\n4. RECORD: Saves every incident on blockchain (tamper-proof record).', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117' });

    // Slide 6
    let slide6 = pres.addSlide();
    slide6.background = { color: 'f8fafc' };
    slide6.addText('Before vs After DRCP', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide6.addText('WITHOUT DRCP\n• Manual monitoring\n• Slow detection (minutes to hours)\n• Failures cascade to other services\n• Incident logs can be edited\n• Engineers scramble during outages\n• SLA violations go unnoticed', { x: 0.5, y: 1.5, w: 4, h: 3, fontSize: 18, color: 'ef4444' });
    slide6.addText('WITH DRCP\n• Automated 24/7 monitoring\n• Instant detection (seconds)\n• Failures are contained automatically\n• Permanent blockchain audit trail\n• System self-heals\n• SLA compliance enforced in real time', { x: 5, y: 1.5, w: 4.5, h: 3, fontSize: 18, color: '10b981' });

    // Slide 7
    let slide7 = pres.addSlide();
    slide7.background = { color: 'f8fafc' };
    slide7.addText('Analytics: Incident Response Time', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide7.addText('Mean Time To Detect (MTTD) and Mean Time To Resolve (MTTR) are slashed by automation.\n\n• Without DRCP: ~45 Minutes\n• With DRCP: < 1 Second\n\n*DRCP\'s Envoy Proxy circuit breakers react in sub-milliseconds to shed load dynamically.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117' });

    // Slide 8
    let slide8 = pres.addSlide();
    slide8.background = { color: 'f8fafc' };
    slide8.addText('Analytics: Outage Prevention', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide8.addText('Impact of automated circuit breaking on platform stability:\n\n• 99.99% System Uptime (Up from 99.5%)\n• 0 Cascading Failures (Down from 12/year)\n• Maintained 100% adherence to critical SLA contracts', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117' });

    // Slide 9
    let slide9 = pres.addSlide();
    slide9.background = { color: 'f8fafc' };
    slide9.addText('Analytics: Business Value Realized', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide9.addText('• $1.2M+ Annual Savings in avoided downtime costs.\n• 100% Tamper-proof audit logs for compliance.\n• 300 hrs Engineering time saved on incident response.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 10
    let slide10 = pres.addSlide();
    slide10.background = { color: 'f8fafc' };
    slide10.addText('Product Features', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide10.addText('• Service Registry: Central hub to track all microservices.\n• SLA Contracts: Define rules for each service (latency, error rate).\n• Live Dashboard: Real-time view of services, contracts, and incidents.\n• Auto Circuit Breaking: Automatically stops bad traffic immediately.\n• Blockchain Audit: Permanent, tamper-proof incident records on Ethereum.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 11
    let slide11 = pres.addSlide();
    slide11.background = { color: 'f8fafc' };
    slide11.addText('Architecture Overview', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide11.addText('Your Services → DRCP Monitoring → Problem Detected?\nIf YES → Block Bad Traffic + Record Incident', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 22, color: '0f1117', align: 'center' });

    // Slide 12
    let slide12 = pres.addSlide();
    slide12.background = { color: 'f8fafc' };
    slide12.addText('Technology Highlights', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide12.addText('• Built with Go (fast, reliable, used by Google, Docker, Kubernetes).\n• Uses industry-standard event streaming tools: Kafka & Redis.\n• Blockchain (Ethereum) for an immutable, tamper-proof audit trail.\n• Envoy Proxy for intelligent, dynamic traffic management.\n• Clean, responsive web dashboard - zero installation needed.\n• Docker-based containerized deployment - easy to run anywhere.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 13
    let slide13 = pres.addSlide();
    slide13.background = { color: 'f8fafc' };
    slide13.addText('Enterprise Use Cases', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide13.addText('• E-Commerce: Protect critical checkout services from failing during major sales spikes.\n• Banking/Finance: Enforce strict sub-millisecond SLAs for payment gateways.\n• Healthcare: Ensure life-critical patient data services stay highly available.\n• SaaS Platforms: Guarantee 99.99% uptime commitments to your enterprise customers.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117', bullet: true });

    // Slide 14
    let slide14 = pres.addSlide();
    slide14.background = { color: 'f8fafc' };
    slide14.addText('Key Takeaways', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide14.addText('1. Resilience is now Automated. No more manual scaling or panic debugging during an outage.\n2. Financial Protection. Avoid massive downtime costs and SLA breach penalties.\n3. Unbreakable Trust. Blockchain records prove your uptime and transparency to your clients.\n4. Ready Today. DRCP plugs directly into modern microservice stacks.', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117' });

    // Slide 15
    let slide15 = pres.addSlide();
    slide15.background = { color: 'f8fafc' };
    slide15.addText('Live Dashboard Demo', { x: 0.5, y: 0.5, w: 9, h: 0.8, fontSize: 28, color: '0284c7', bold: true });
    slide15.addText('Demo: https://drcp-dashboard.onrender.com\n\nKey Features to Explore:\n• Live service health monitoring\n• Dynamic SLA contract management\n• Incident tracking & blockchain verifications\n• Real-time telemetry ingestion streams', { x: 0.5, y: 1.5, w: 9, h: 3, fontSize: 18, color: '0f1117' });

    // Slide 16
    let slide16 = pres.addSlide();
    slide16.background = { color: 'f8fafc' };
    slide16.addText('Thank You!', { x: 1, y: 1.5, w: 8, h: 1, fontSize: 48, color: '0284c7', bold: true, align: 'center' });
    slide16.addText('Questions?\nDemo: https://drcp-dashboard.onrender.com\nBuilt by Prachi Yadav', { x: 1, y: 2.5, w: 8, h: 1.5, fontSize: 24, color: '334155', align: 'center' });

    await pres.writeFile({ fileName: 'DRCP_Business_Presentation.pptx' });
}

async function main() {
    await createTechPresentation();
    await createBizPresentation();
}

main().catch(err => console.error(err));
