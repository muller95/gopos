package com.mycompany.app;

/**
 * Hello world!
 *
 */

import org.xhtmlrenderer.pdf.ITextRenderer;
import org.xhtmlrenderer.simple.*;

import java.awt.print.Pageable;
import java.awt.print.PrinterJob;

  
import java.io.*;
import java.util.*;
import java.nio.file.*;
import javax.print.*;
import javax.print.attribute.*;
import javax.print.attribute.standard.*;

enum PrinterType { Check,  Order, Both }

class Printer {
    private String name;
    private PrinterType type;

    public Printer(String name, PrinterType type) {
        this.name = name;
        this.type = type;
    }
    public synchronized void print(File file) {
        try {
            XHTMLPanel panel = new XHTMLPanel();
            panel.setDocument(file);
            XHTMLPrintable printable = new XHTMLPrintable(panel);

            
            AttributeSet attrSet = new HashAttributeSet();
            attrSet.add(new PrinterName(name, null));
            PrintService[] printServices =  PrintServiceLookup.lookupPrintServices(null, attrSet);
            System.out.println("---------------SUPPORTED FLAVOURS------------------");
            DocFlavor[] flavours = printServices[0].getSupportedDocFlavors();
            for(int i = 0; i < flavours.length; i++)
                System.out.println(flavours[i]);
            System.out.println("---------------------------------------------------");   
            
            PrinterJob job = PrinterJob.getPrinterJob();
            job.setPrintService(printServices[0]);
            job.setPrintable(printable);            
            job.print();

            file.delete();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public PrinterType getType() {
        return type;
    }
}
class ToPdfRunnable implements Runnable {
        private String data;
        private String path;
        public ToPdfRunnable(String data, String path) {
            this.data = data;
            this.path = path;
        }   

        public void run() {        
            try {              
                FileOutputStream os = new FileOutputStream(path);
                ITextRenderer renderer = new ITextRenderer();
                renderer.setDocumentFromString(data);
                renderer.layout();
                renderer.createPDF(os);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
}

class PrintRunnable implements Runnable {
        private File file;
        private Printer printer;
        public PrintRunnable(File file, Printer printer) {
            this.file = file;
            this.printer = printer;
        }   

        public void run() {        
            try {  
                printer.print(file);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
}


public class App {
    
    public static void main( String[] args ) {
        HashMap<String, PrinterType> printerNames = new HashMap<String, PrinterType>();
        Map<String, String> env = System.getenv();
        String goposCheckPath = env.get("GOPOS_CHECK_PATH");
        String goposOrderPath = env.get("GOPOS_ORDER_PATH");
        String goposPrintserviceCheckTmpPath = env.get("GOPOS_PRINTSERVICE_CHECK_TMP_PATH");
        String goposPrintserviceOrderTmpPath = env.get("GOPOS_PRINTSERVICE_ORDER_TMP_PATH");    
        String goposPrintserviceCheckPrinterNames = env.get("GOPOS_PRINTSERVICE_CHECK_" + 
            "PRINTER_NAME");
        String goposPrintserviceOrderPrinterNames = env.get("GOPOS_PRINTSERVICE_ORDER_" + 
            "PRINTER_NAME");

        if (goposCheckPath == null) {
            System.err.println("GOPOS_CHECK_PATH is not set.");
            return;
        }

        if (goposOrderPath == null) {
            System.err.println("GOPOS_ORDER_PATH is not set.");
            return;
        }

        if (goposPrintserviceCheckTmpPath == null) {
            System.err.println("GOPOS_PRINTSERVICE_CHECK_TMP_PATH is not set.");
            return;
        }

        if (goposPrintserviceOrderTmpPath == null) {
            System.err.println("GOPOS_PRINTSERVICE_ORDER_TMP_PATH is not set.");
            return;
        }
        
        if (goposPrintserviceCheckPrinterNames == null) {
            System.err.println("GOPOS_PRINTSERVICE_CHECK_PRINTER_NAME is not set");
            return;
        }

        if (goposPrintserviceOrderPrinterNames == null) {
            System.err.println("GOPOS_PRINTSERVICE_ORDER_PRINTER_NAME is not set");
            return;
        }

        System.out.println("---------LIST ALL PRINTERS--------------------");
        PrintService[] printServices = PrintServiceLookup.lookupPrintServices(null, null);
        System.out.println("Number of print services: " + printServices.length);

        for (PrintService printer : printServices)
            System.out.println("Printer: " + printer.getName()); 
        System.out.println("---------------------------------------");

        System.out.println("printer check names: " + goposPrintserviceCheckPrinterNames);
        System.out.println("printer order names: " + goposPrintserviceOrderPrinterNames);
        System.out.println("check path: " + goposCheckPath); 
        System.out.println("order path: " + goposOrderPath); 
        System.out.println("tmp check path: " + goposPrintserviceCheckTmpPath); 
        System.out.println("tmp order path: " + goposPrintserviceOrderTmpPath);

        String checkPrinterNames[] = goposPrintserviceCheckPrinterNames.split(":");
        for (int i = 0; i < checkPrinterNames.length; i++)
            printerNames.put(checkPrinterNames[i].trim(), PrinterType.Check);

        String orderPrinterNames[] = goposPrintserviceOrderPrinterNames.split(":");
        for (int i = 0; i < orderPrinterNames.length; i++) {
            String key = orderPrinterNames[i].trim();
            if (printerNames.containsKey(key)) 
                printerNames.put(key, PrinterType.Both);
            else 
                printerNames.put(key, PrinterType.Order);                
        }

        ArrayList<Printer> printers = new ArrayList<Printer>();
        for(Map.Entry<String, PrinterType> entry : printerNames.entrySet())
            printers.add(new Printer(entry.getKey(), entry.getValue()));


        try {
            File checkDir = new File(goposCheckPath);
            File orderDir = new File(goposOrderPath);
            File checkPdfDir = new File(goposPrintserviceCheckTmpPath);
            File orderPdfDir = new File(goposPrintserviceOrderTmpPath);

            for (;;) {
                Thread.sleep(1000);
                File checks[] = checkDir.listFiles();   
                File orders[] = orderDir.listFiles();
                File checkPdfs[] = checkPdfDir.listFiles();   
                File orderPdfs[] = orderPdfDir.listFiles();
                
                // System.out.println("number of checks: " + Integer.toString(checks.length));                
                /*for (int i = 0; i < checks.length; i++) {
                    BufferedReader reader = new BufferedReader(new FileReader(checks[i]));
                    String data = "", tmp = "";
                    while ((tmp = reader.readLine()) != null)
                        data += "\n" + tmp;

                    String path = goposPrintserviceCheckTmpPath + checks[i].getName().
                        substring(0, checks[i].getName().length() - 4) + "pdf";
                    
                    new Thread(new ToPdfRunnable(data, path)).start();
                    reader.close();
                    checks[i].delete();
                }*/

                // System.out.println("number of orders: " + Integer.toString(orders.length));                
                for (int i = 0; i < orders.length; i++) {
                    String path = goposPrintserviceOrderTmpPath + orders[i].getName();
                    Files.copy(Paths.get(orders[i].getPath()), Paths.get(path), 
                        StandardCopyOption.REPLACE_EXISTING);
                    orders[i].delete();

                    
                    for (int j = 0; j < printers.size(); j++) {
                        PrinterType currType = printers.get(j).getType(); 
                        if (currType == PrinterType.Order || currType == PrinterType.Both)
                            new Thread(new PrintRunnable(new File(path), printers.get(j))).start();
                    }
                }

                /*for (int i = 0; i < checkPdfs.length; i++) {
                    byte pdfData[] = Files.readAllBytes(Paths.get(checkPdfs[i].getPath()));
                    System.out.println("print check:" + checkPdfs[i].getPath());
                    checkPdfs[i].delete();
                    for (int j = 0; j < printers.size(); j++) {
                        PrinterType currType = printers.get(j).getType(); 
                        if (currType == PrinterType.Check || currType == PrinterType.Both)
                            new Thread(new PrintRunnable(pdfData, printers.get(j))).start();
                    }
                }

                for (int i = 0; i < orderPdfs.length; i++) {
                    byte pdfData[] = Files.readAllBytes(Paths.get(orderPdfs[i].getPath()));
                    System.out.println("print order:" + orderPdfs[i].getPath());
                    orderPdfs[i].delete();
                    for (int j = 0; j < printers.size(); j++) {
                        PrinterType currType = printers.get(j).getType(); 
                        if (currType == PrinterType.Order || currType == PrinterType.Both)
                            new Thread(new PrintRunnable(pdfData, printers.get(j))).start();
                    }
                }*/
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
