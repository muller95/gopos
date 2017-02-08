package com.mycompany.app;

/**
 * Hello world!
 *
 */

import org.xhtmlrenderer.pdf.ITextRenderer;
import java.io.*;

public class App 
{
    public static void main( String[] args )
    {
        String HTML = "/home/vadim/Sources/gopos/server/5_28.html";
        String PDF = "/home/vadim/test.pdf";

        try {
            String url = new File(HTML).toURI().toURL().toString();
            FileOutputStream os = new FileOutputStream(PDF);

            ITextRenderer renderer = new ITextRenderer();
            renderer.setDocument(url);
            renderer.layout();
            renderer.createPDF(os);

            os.close();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
