use anyhow::Result;
use clap::Parser;

mod cli;
mod capture;
mod config;
mod detection;
mod error;
mod firewall;
mod stats;
mod storage;

use cli::{Cli, Commands};

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::init();
    
    let cli = Cli::parse();
    
    match cli.command {
        Commands::ListInterfaces => {
            println!("🌐 Available Network Interfaces:");
            println!("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
            
            let interfaces = capture::list_interfaces()?;
            
            for (idx, iface) in interfaces.iter().enumerate() {
                println!("\n[{}] {}", idx + 1, iface.name);
                println!("    Description: {}", iface.description.as_deref().unwrap_or("N/A"));
                println!("    MAC: {}", iface.mac.as_deref().unwrap_or("N/A"));
                
                if !iface.ips.is_empty() {
                    println!("    IPs:");
                    for ip in &iface.ips {
                        println!("      - {}", ip);
                    }
                }
            }
        }
        
        Commands::Monitor {
            interface,
            db_path,
            config_file,
            verbose,
        } => {
            use colored::Colorize;
            
            println!("{}", "🛡️  NetGuard - Network Monitor".bright_cyan().bold());
            println!("{}", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━".bright_black());
            
            // Load configuration
            let config = if let Some(config_path) = config_file {
                config::Config::from_file(&config_path)?
            } else {
                config::Config::default()
            };
            
            // Initialize storage if database path provided
            let storage = if let Some(db) = db_path {
                Some(storage::Storage::new(&db)?)
            } else {
                None
            };
            
            // Start monitoring
            let monitor = capture::Monitor::new(interface, config, storage, verbose)?;
            
            println!("\n{}", "Starting packet capture...".green());
            println!("{}", "Press Ctrl+C to stop".yellow());
            println!();
            
            monitor.start().await?;
        }
        
        Commands::Stats {
            interface,
            db_path,
            history,
        } => {
            if history {
                if let Some(db) = db_path {
                    let storage = storage::Storage::new(&db)?;
                    stats::display_historical_stats(&storage)?;
                } else {
                    eprintln!("Error: --db required for historical stats");
                }
            } else {
                let stats_monitor = stats::StatsMonitor::new(interface)?;
                stats_monitor.display_realtime().await?;
            }
        }
        
        Commands::Rules { subcommand } => {
            use cli::RulesCommands;
            
            match subcommand {
                RulesCommands::List => {
                    firewall::list_rules()?;
                }
                RulesCommands::Add {
                    block,
                    ip,
                    port,
                    protocol,
                } => {
                    firewall::add_rule(block, ip, port, protocol)?;
                }
                RulesCommands::Remove { id } => {
                    firewall::remove_rule(id)?;
                }
                RulesCommands::Load { file } => {
                    firewall::load_rules(&file)?;
                }
            }
        }
        
        Commands::Alerts {
            db_path,
            severity,
            export,
            limit,
        } => {
            let storage = storage::Storage::new(&db_path)?;
            storage::display_alerts(&storage, severity, export, limit)?;
        }
    }
    
    Ok(())
}