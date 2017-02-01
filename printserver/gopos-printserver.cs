using System;

namespace GoposPrintserver {
    class Program {
        static string GoposPrintserverIP;
        static string GoposPrintserverPort;
        static string GoposPrintserverKitchenPrinterName;
        static string GoposPrintserverCheckPrinterName;

        static void Main(string[] args) {
            GoposPrintserverIP = Environment.GetEnvironmentVariable("GOPOS_PRINTSERVER_IP");
            GoposPrintserverPort = Environment.GetEnvironmentVariable("GOPOS_PRINTSERVER_PORT");
            GoposPrintserverKitchenPrinterName = Environment.GetEnvironmentVariable("GOPOS_" + 
                "PRINTSERVER_KITCHEN_PRINTER_NAME");
            GoposPrintserverCheckPrinterName = Environment.GetEnvironmentVariable("GOPOS_" + 
                "PRINTSERVER_CHECK_PRINTER_NAME");

            if (GoposPrintserverIP == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_IP is not set");
                Environment.Exit(0);
            }

            if (GoposPrintserverPort == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_PORT is not set");
                Environment.Exit(0);
            }

            if (GoposPrintserverKitchenPrinterName == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_KITCHEN_PRINTER_NAME is not set");
                Environment.Exit(0);
            }

            if (GoposPrintserverCheckPrinterName == null) {
                Console.WriteLine("GOPOS_PRINTSERVER_KITCHEN_PRINTER_NAME is not set");
                Environment.Exit(0);
            }

            Console.WriteLine("here");
        }
    }    
}