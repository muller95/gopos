using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;

namespace GoposPrintserver {
    class Program {
        enum PrinterType { Check, Order, Both }
        class Printer {
            private string name;
            private Queue<string[]> packages;
            private PrinterType type;

            public Printer(string name, PrinterType type) {
                this.name = name;
                this.type = type;   
                packages = new Queue<String[]>();
            }
        }
        
        static List<Printer> printers;

        static void Main(string[] args) {
            string goposPrintserverPort;
            string goposPrintserverOrderPrinterNames, goposPrintserverCheckPrinterNames;
            string[] orderPrinterNames, checkPrinterNames;
            Dictionary<string, PrinterType> printerDictionary;
            TcpListener tcpListener;

            goposPrintserverPort = Environment.GetEnvironmentVariable("GOPOS_PRINTSERVER_PORT");
            goposPrintserverOrderPrinterNames = Environment.GetEnvironmentVariable("GOPOS_" + 
                "PRINTSERVER_ORDER_PRINTER_NAME");
            goposPrintserverCheckPrinterNames = Environment.GetEnvironmentVariable("GOPOS_" + 
                "PRINTSERVER_CHECK_PRINTER_NAME");

            if (goposPrintserverPort == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_PORT is not set");
                Environment.Exit(0);
            }

            if (goposPrintserverOrderPrinterNames == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_KITCHEN_PRINTER_NAME is not set");
                Environment.Exit(0);
            }

            if (goposPrintserverCheckPrinterNames == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_KITCHEN_PRINTER_NAME is not set");
                Environment.Exit(0);
            }

            orderPrinterNames = goposPrintserverOrderPrinterNames.Split(new char[] { ':' }, 
                StringSplitOptions.None);
            checkPrinterNames = goposPrintserverCheckPrinterNames.Split(new char[] { ':' }, 
                StringSplitOptions.None);

            printerDictionary = new Dictionary<string, PrinterType>();
            for (int i = 0; i < orderPrinterNames.Length; i++)
                printerDictionary.Add(orderPrinterNames[i].Trim(new char[] { ' ' }), 
                    PrinterType.Order);
            for (int i = 0; i < checkPrinterNames.Length; i++) {
                string key;

                key = orderPrinterNames[i].Trim(new char[] { ' ' });
                if (printerDictionary.ContainsKey(key))
                    printerDictionary[key] = PrinterType.Both;
                else
                    printerDictionary.Add(key, PrinterType.Check);
            }

            printers = new List<Printer>();
            foreach(KeyValuePair<string, PrinterType> entry in printerDictionary)
                printers.Add(new Printer(entry.Key, entry.Value));
            
            tcpListener = new TcpListener(IPAddress.Any, Convert.ToInt32(goposPrintserverPort));
            tcpListener.Start();
            while (true) {
                TcpClient tcpClient;

                
            }
        }
    }    
}