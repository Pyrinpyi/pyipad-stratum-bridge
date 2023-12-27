# Hive configuration
## Local bridge (1 bridge per miner)
This assumes that you will be running the ks_bridge on the mining computer. I've provided 2 template flight sheets that can be imported to jumpstart the config process 

---

I'll be using bzminer here, so I'm focusing on the file [here](misc/hive-templates/pyi_srb_local_adapter.json). Hive gives the option to import from a file or from your clipboard; I'll be copying directly from the clipboard in this example. 

Navigate to the template that you want and hit the 'copy file contents' button
![image](https://github.com/Pyrinpyi/pyipad-stratum-bridge/assets/138370658/a72add34-0629-42dd-a3fc-4ff10cc7c89b)


This will load the contents of the file in your clipboard. Now go to hive, create a new flight sheet, and click the import from clipboard button


Once the import happens you'll be presented with a flightsheet filled with angry little red boxes like below.
![image](https://github.com/Pyrinpyi/pyipad-stratum-bridge/assets/138370658/acd36eb5-eb2e-450d-81a8-b4e25fde005a)



 Set the coin & wallet for both miners in the flight sheet, and choose 'Configure in miner' for the pool. Strangly you cannot edit the miner configs until you touch the pool config, so re-set it to 'Configure in miner' if the UI is being weird. 
 
At this point the flight sheet should be happy:
![image](https://github.com/Pyrinpyi/pyipad-stratum-bridge/assets/138370658/e9019889-dfb0-4378-91e0-513832896c65)


Last step, you just need to edit the custom miner section with the address of the pyrin node you're mining to. This needs the full ip/port so it'll look similar to `192.168.0.12:13110` save the changes after you're finished

![image](https://github.com/Pyrinpyi/pyipad-stratum-bridge/assets/138370658/6756d318-76a2-4eec-b784-15db3dc6d02a)



You should be good to create the flight sheet now. However, you should verify the configurations in both miners before you run it. Never run something blindly off the internet :)

This is what the configuration generally should look like
![image](https://github.com/Pyrinpyi/pyipad-stratum-bridge/assets/138370658/60092d7f-17c3-4ccb-b986-5c37e64249d5)



The important parts are the 'Install URL' (make sure that points back to this repo somewhere) and double check that your wallet address is set correctly. 

You can verify that the bridge is working by either feching miner logs, or checking on the main miner--if it's mining then the bridge is working!



## Remote bridge (1 bridge for N miners)
This is the preferred configuration as it's a bit easier to setup and the bridge can easily handle multiple miners.

Ideally you'd run the bridge on the same machine that your pyipad node is running on. You can follow any of the installation methods defined in the main readme of the repo.

Once the bridge is running, miners can be configured simply by replacing the pool URL with the ip:port of the bridge. In the example below the bridge is running on `192.168.0.25:5555`

![image](https://github.com/Pyrinpyi/pyipad-stratum-bridge/assets/138370658/28c62e9e-47ed-418c-9558-dad54d417688)





